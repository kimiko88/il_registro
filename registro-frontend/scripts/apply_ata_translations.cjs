const fs = require('fs');
const path = require('path');
const vm = require('vm');

const localeMap = {
  'it-IT': './translations/it.cjs',
  'en-US': './translations/en.cjs',
  'de-DE': './translations/de.cjs',
  'es-ES': './translations/es.cjs',
  'fr-FR': './translations/fr.cjs',
  'ro-RO': './translations/ro.cjs',
  'ru-RU': './translations/ru.cjs',
  'uk-UA': './translations/uk.cjs',
  'sq-AL': './translations/sq.cjs',
  'ar-SA': './translations/ar.cjs',
  'zh-CN': './translations/zh.cjs'
};

function deepMerge(target, source) {
  for (const key of Object.keys(source)) {
    if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
      if (!target[key] || typeof target[key] !== 'object' || Array.isArray(target[key])) {
        target[key] = {};
      }
      deepMerge(target[key], source[key]);
    } else {
      target[key] = source[key];
    }
  }
  return target;
}

Object.entries(localeMap).forEach(([loc, transRelPath]) => {
  const transPath = path.resolve(__dirname, transRelPath);
  const trans = require(transPath);

  const targetFile = path.resolve(__dirname, '../src/i18n', loc, 'index.js');
  if (!fs.existsSync(targetFile)) {
    console.error(`Missing target file for ${loc}: ${targetFile}`);
    return;
  }

  const rawCode = fs.readFileSync(targetFile, 'utf8');
  const cjsCode = rawCode.replace(/^\s*export\s+default\s+/m, 'module.exports = ');
  const m = { exports: {} };
  const script = new vm.Script(cjsCode);
  const context = vm.createContext({ module: m, exports: m.exports, console });
  script.runInContext(context);
  const data = m.exports;

  // Ensure sections exist
  data.onboarding = data.onboarding || {};
  data.onboardingExtra = data.onboardingExtra || {};
  data.help = data.help || {};
  data.guideCenter = data.guideCenter || {};

  // Deep merge ATA translations into each section
  deepMerge(data.onboarding, trans.onboarding);
  deepMerge(data.onboardingExtra, trans.onboardingExtra);
  deepMerge(data.help, trans.help);
  deepMerge(data.guideCenter, trans.guideCenter);

  // Write back formatted ES module
  const output = `export default ${JSON.stringify(data, null, 2)}\n`;
  fs.writeFileSync(targetFile, output, 'utf8');
  console.log(`[OK] Updated ${loc} (${output.length} bytes, ~${output.split('\n').length} lines)`);
});
