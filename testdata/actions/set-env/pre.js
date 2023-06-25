const fs = require('fs');

const b = fs.readFileSync(process.env.GITHUB_ENV);
fs.writeFileSync(process.env.GITHUB_ENV, `${b}GENERAL=kenobi\n`);
