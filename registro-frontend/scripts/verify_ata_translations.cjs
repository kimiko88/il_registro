const fs = require('fs');
const path = require('path');
const vm = require('vm');

const locales = ['it-IT', 'en-US', 'de-DE', 'es-ES', 'fr-FR', 'ro-RO', 'ru-RU', 'uk-UA', 'sq-AL', 'ar-SA', 'zh-CN'];
const roles = ['assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico', 'dsga'];

const roleGuides = {
  assistente_amministrativo: ['dashboard', 'personnel_desk', 'sidi', 'attendance', 'verbali'],
  collaboratore_ds: ['dashboard', 'substitutions', 'emergency', 'strike', 'verbali'],
  collaboratore_scolastico: ['dashboard', 'visitors', 'early_exits', 'maintenance', 'badge'],
  dsga: ['dashboard', 'timecard', 'personnel_desk', 'sidi', 'strike']
};

let allPassed = true;

locales.forEach(loc => {
  const filePath = path.resolve(__dirname, '../src/i18n', loc, 'index.js');
  const code = fs.readFileSync(filePath, 'utf8');
  const cjsCode = code.replace(/^\s*export\s+default\s+/m, 'module.exports = ');
  const m = { exports: {} };
  const script = new vm.Script(cjsCode);
  const context = vm.createContext({ module: m, exports: m.exports, console });
  script.runInContext(context);
  const data = m.exports;

  const errors = [];

  roles.forEach(role => {
    // 1. onboarding
    for (let s = 1; s <= 5; s++) {
      if (!data.onboarding?.[role]?.[`step${s}_title`]) {
        errors.push(`missing onboarding.${role}.step${s}_title`);
      }
      if (!data.onboarding?.[role]?.[`step${s}_desc`]) {
        errors.push(`missing onboarding.${role}.step${s}_desc`);
      }
    }

    // 2. onboardingExtra
    for (let s = 1; s <= 5; s++) {
      const b = data.onboardingExtra?.[role]?.[`step${s}_bullets`];
      if (!Array.isArray(b) || b.length < 3) {
        errors.push(`missing or invalid onboardingExtra.${role}.step${s}_bullets`);
      }
    }

    // 3. help FAQs
    for (let q = 1; q <= 5; q++) {
      if (!data.help?.[role]?.[`q${q}`]) {
        errors.push(`missing help.${role}.q${q}`);
      }
      if (!data.help?.[role]?.[`a${q}`]) {
        errors.push(`missing help.${role}.a${q}`);
      }
    }

    // 4. guideCenter
    const guides = roleGuides[role];
    guides.forEach(g => {
      const guide = data.guideCenter?.[role]?.[g];
      if (!guide) {
        errors.push(`missing guideCenter.${role}.${g}`);
      } else {
        ['title', 'desc', 'content', 'categoryLabel', 'categoryDesc'].forEach(prop => {
          if (!guide[prop]) {
            errors.push(`missing guideCenter.${role}.${g}.${prop}`);
          }
        });
      }
    });
  });

  if (errors.length > 0) {
    allPassed = false;
    console.error(`[FAIL] ${loc}: ${errors.length} errors:`, errors.slice(0, 5));
  } else {
    console.log(`[PASS] ${loc}: all ATA tour & guide keys verified successfully!`);
  }
});

if (allPassed) {
  console.log('\nSUCCESS: All 11 languages have complete, verified ATA translations!');
  process.exit(0);
} else {
  console.error('\nFAILURE: Some translations are missing or invalid.');
  process.exit(1);
}
