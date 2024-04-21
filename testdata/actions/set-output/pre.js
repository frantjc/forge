const fs = require('fs');

console.log('::set-output name=hello::there');
const b = fs.readFileSync(process.env.GITHUB_OUTPUT);
fs.writeFileSync(process.env.GITHUB_OUTPUT, `${b}general=kenobi\n`);
