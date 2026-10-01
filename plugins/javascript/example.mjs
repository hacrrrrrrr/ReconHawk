const target = process.argv[2] ?? "";
console.log(JSON.stringify({ plugin: "javascript-example", target, type: "enrichment" }));
