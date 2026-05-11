const cp = require("child_process");

process.on("message", msg => {
    if (msg.command === "run") {
        const exe = msg.program;
        const args = msg.args || [];

        const proc = cp.spawn(exe, args, { stdio: "inherit" });

        proc.on("exit", code => {
            process.exit(code);
        });
    }
});
