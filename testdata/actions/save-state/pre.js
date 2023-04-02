const fs = require('fs');

console.log('::save-state name=hello::there');
const b = fs.readFileSync(process.env.GITHUB_STATE);
fs.writeFileSync(process.env.GITHUB_STATE, `${b}general=kenobi\n`);
