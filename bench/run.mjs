import fs from "node:fs";
import path from "node:path";
import os from "node:os";
import { spawn, spawnSync, execSync } from "node:child_process";
import crypto from "node:crypto";
import readline from "node:readline";

const getRssBytes = (pid) => {
    try {
        if (os.platform() === 'linux') {
            const status = fs.readFileSync(`/proc/${pid}/status`, 'utf8');
            const match = status.match(/VmRSS:\s+(\d+)\s+kB/);
            if (match) return parseInt(match[1], 10) * 1024;
        } else if (os.platform() === 'win32') {
            const output = execSync(`powershell -Command "(Get-Process -Id ${pid}).WorkingSet64"`, { encoding: 'utf8' });
            const val = parseInt(output.trim(), 10);
            if (!isNaN(val)) return val;
        }
    } catch (e) {
        // ignore
    }
    return null;
};

const percentile99 = (arr) => {
    if (!arr || arr.length === 0) return 0;
    const sorted = [...arr].sort((a, b) => Number(a) - Number(b));
    const index = Math.ceil(0.99 * sorted.length) - 1;
    return sorted[index];
};

const runSingleRequest = async (command, args, workloadArray) => {
    return new Promise((resolve, reject) => {
        const start = process.hrtime.bigint();
        const proc = spawn(command, args, { stdio: ['pipe', 'pipe', 'ignore'] });

        let returned = false;
        const rl = readline.createInterface({ input: proc.stdout, crlfDelay: Infinity });

        proc.on('error', reject);

        rl.on('line', (line) => {
            if (returned) return;
            returned = true;
            const end = process.hrtime.bigint();
            const elapsed = Number(end - start) / 1e6;
            proc.kill('SIGTERM');
            resolve({ elapsed, rss: getRssBytes(proc.pid) });
        });

        // send first request to trigger processing
        proc.stdin.write(JSON.stringify(workloadArray[0]) + '\n');
    });
};

const runPersistent = async (command, args, workloadArray, numRequests) => {
    return new Promise((resolve, reject) => {
        const proc = spawn(command, args, { stdio: ['pipe', 'pipe', 'ignore'] });
        const rl = readline.createInterface({ input: proc.stdout, crlfDelay: Infinity });

        let latencies = [];
        let i = 0;
        let startReq = 0n;
        let maxRss = null;

        const overallStart = process.hrtime.bigint();

        proc.on('error', reject);

        const sendNext = () => {
            if (i >= numRequests) {
                const overallEnd = process.hrtime.bigint();
                const totalElapsedSeconds = Number(overallEnd - overallStart) / 1e9;
                proc.kill('SIGTERM');
                resolve({
                    latencies,
                    throughput: numRequests / totalElapsedSeconds,
                    maxRss
                });
                return;
            }
            startReq = process.hrtime.bigint();
            proc.stdin.write(JSON.stringify(workloadArray[i % workloadArray.length]) + '\n');
            const rss = getRssBytes(proc.pid);
            if (rss !== null) {
                if (maxRss === null || rss > maxRss) maxRss = rss;
            }
        };

        rl.on('line', (line) => {
            const endReq = process.hrtime.bigint();
            latencies.push(Number(endReq - startReq) / 1e6);
            i++;
            sendNext();
        });

        sendNext();
    });
};

const main = async () => {
    const args = process.argv.slice(2);
    const isQuick = args.includes('--quick');
    let outputPath = 'bench/results.json';
    const outIdx = args.indexOf('--output');
    if (outIdx !== -1 && outIdx + 1 < args.length) {
        outputPath = args[outIdx + 1];
    }

    const coldStarts = isQuick ? 3 : 20;
    const persistentRequests = isQuick ? 30 : 1000;

    // build go binary if needed
    if (!fs.existsSync('bin/tinycolor') && !fs.existsSync('bin/tinycolor.exe')) {
        spawnSync('go', ['-C', 'src', 'build', '-o', '../bin/tinycolor', './cmd/tinycolor-compat']);
    }
    const goBinary = fs.existsSync('bin/tinycolor.exe') ? 'bin/tinycolor.exe' : 'bin/tinycolor';

    const workloadContent = fs.readFileSync('bench/workload.jsonl', 'utf8');
    const workloadLines = workloadContent.trim().split('\n');
    const workloadArray = workloadLines.map(l => JSON.parse(l));
    const sha256 = crypto.createHash('sha256').update(workloadContent).digest('hex');

    const implementations = {};

    for (const impl of ['javascript', 'go']) {
        let command, cmdArgs;
        if (impl === 'javascript') {
            command = 'node';
            cmdArgs = ['compat/js-runner.mjs'];
        } else {
            command = path.join(process.cwd(), goBinary);
            cmdArgs = [];
        }

        let startupLatencies = [];
        let maxRss = null;

        for (let i = 0; i < coldStarts; i++) {
            const res = await runSingleRequest(command, cmdArgs, workloadArray);
            startupLatencies.push(res.elapsed);
            if (res.rss !== null) {
                if (maxRss === null || res.rss > maxRss) maxRss = res.rss;
            }
        }

        const persistentRes = await runPersistent(command, cmdArgs, workloadArray, persistentRequests);
        if (persistentRes.maxRss !== null) {
            if (maxRss === null || persistentRes.maxRss > maxRss) maxRss = persistentRes.maxRss;
        }

        implementations[impl] = {
            startupP99Ms: percentile99(startupLatencies),
            latencyP99Ms: percentile99(persistentRes.latencies),
            throughputOpsPerSecond: persistentRes.throughput,
            peakRssBytes: maxRss
        };
    }

    const nodeVersion = process.version;
    const goVersionMatch = spawnSync('go', ['version'], { encoding: 'utf8' }).stdout.match(/go version (\S+)/);
    const goVersion = goVersionMatch ? goVersionMatch[1] : 'unknown';

    const result = {
        generatedAt: new Date().toISOString(),
        os: os.type(),
        architecture: os.arch(),
        cpu: os.cpus()[0].model,
        versions: {
            node: nodeVersion,
            go: goVersion
        },
        workloadSha256: sha256,
        samples: {
            coldStarts,
            requests: persistentRequests
        },
        implementations
    };

    fs.writeFileSync(outputPath, JSON.stringify(result, null, 2));
};

main().catch(e => {
    console.error(e);
    process.exit(1);
});
