#!/usr/bin/env node
/**
 * scripts/bump-version.mjs
 *
 * Automated, atomic semantic version manager for "il_registro" (Registrov2).
 * Synchronizes version across:
 *  - VERSION (root single source of truth)
 *  - registro-frontend/package.json
 *  - registro-backend/pkg/version/version.go
 *  - publiccode.yml (softwareVersion & releaseDate)
 *  - CHANGELOG.md (adds/updates section header)
 *
 * Usage:
 *   node scripts/bump-version.mjs [patch | minor | major | beta | <custom-version>]
 *   node scripts/bump-version.mjs --check
 *   node scripts/bump-version.mjs patch --git
 */

import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { execSync } from 'node:child_process'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const rootDir = path.resolve(__dirname, '..')

// File paths
const files = {
  versionFile: path.join(rootDir, 'VERSION'),
  frontendPkg: path.join(rootDir, 'registro-frontend', 'package.json'),
  backendVersionGo: path.join(rootDir, 'registro-backend', 'pkg', 'version', 'version.go'),
  publiccode: path.join(rootDir, 'publiccode.yml'),
  changelog: path.join(rootDir, 'CHANGELOG.md')
}

function getTodayISO() {
  return new Date().toISOString().split('T')[0]
}

function readCurrentVersion() {
  if (fs.existsSync(files.versionFile)) {
    return fs.readFileSync(files.versionFile, 'utf-8').trim()
  }
  if (fs.existsSync(files.frontendPkg)) {
    const pkg = JSON.parse(fs.readFileSync(files.frontendPkg, 'utf-8'))
    return pkg.version
  }
  return '1.0.0-beta'
}

function parseSemver(v) {
  const clean = v.replace(/^v/, '').trim()
  const match = clean.match(/^(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z.-]+))?$/)
  if (!match) {
    return null
  }
  return {
    major: parseInt(match[1], 10),
    minor: parseInt(match[2], 10),
    patch: parseInt(match[3], 10),
    prerelease: match[4] || null,
    raw: clean
  }
}

function computeNextVersion(currentStr, bumpType) {
  const parsed = parseSemver(currentStr)
  if (!parsed) {
    throw new Error(`Current version "${currentStr}" is not valid semver.`)
  }

  const { major, minor, patch, prerelease } = parsed

  switch (bumpType.toLowerCase()) {
    case 'patch':
      if (prerelease) {
        // e.g. 1.0.0-beta -> 1.0.0
        return `${major}.${minor}.${patch}`
      }
      return `${major}.${minor}.${patch + 1}`

    case 'minor':
      return `${major}.${minor + 1}.0`

    case 'major':
      return `${major + 1}.0.0`

    case 'beta':
    case 'prerelease':
      if (!prerelease) {
        return `${major}.${minor}.${patch + 1}-beta.1`
      }
      const betaMatch = prerelease.match(/beta\.?(\d+)?/)
      if (betaMatch && betaMatch[1]) {
        const nextNum = parseInt(betaMatch[1], 10) + 1
        return `${major}.${minor}.${patch}-beta.${nextNum}`
      }
      return `${major}.${minor}.${patch}-beta.1`

    default:
      // If user supplied explicit version string (e.g. "1.2.0")
      const customParsed = parseSemver(bumpType)
      if (!customParsed) {
        throw new Error(
          `Invalid bump type or version "${bumpType}". Must be one of: patch, minor, major, beta, or a valid semver string like 1.0.1`
        )
      }
      return customParsed.raw
  }
}

function checkVersions() {
  console.log('🔍 Checking version synchronization across repository...\n')
  const rootVer = readCurrentVersion()

  const pkgVer = fs.existsSync(files.frontendPkg)
    ? JSON.parse(fs.readFileSync(files.frontendPkg, 'utf-8')).version
    : 'MISSING'

  let goVer = 'MISSING'
  if (fs.existsSync(files.backendVersionGo)) {
    const goContent = fs.readFileSync(files.backendVersionGo, 'utf-8')
    const match = goContent.match(/Version\s*=\s*"([^"]+)"/)
    if (match) goVer = match[1]
  }

  let publiccodeVer = 'MISSING'
  if (fs.existsSync(files.publiccode)) {
    const pcContent = fs.readFileSync(files.publiccode, 'utf-8')
    const match = pcContent.match(/softwareVersion:\s*([^\r\n]+)/)
    if (match) publiccodeVer = match[1].trim()
  }

  console.log(`  • ROOT VERSION file:               ${rootVer}`)
  console.log(`  • registro-frontend/package.json:  ${pkgVer}`)
  console.log(`  • registro-backend pkg/version:    ${goVer}`)
  console.log(`  • publiccode.yml:                  ${publiccodeVer}`)

  const isSynced = rootVer === pkgVer && rootVer === goVer && rootVer === publiccodeVer
  if (isSynced) {
    console.log('\n✅ All components are perfectly synchronized on version:', rootVer)
  } else {
    console.log('\n⚠️  Versions are OUT OF SYNC! Run `npm run bump [version]` to synchronize.')
    process.exitCode = 1
  }
}

function updateFiles(newVersion) {
  const today = getTodayISO()

  // 1. Root VERSION
  fs.writeFileSync(files.versionFile, `${newVersion}\n`, 'utf-8')
  console.log(`  ✓ Updated ${path.relative(rootDir, files.versionFile)} -> ${newVersion}`)

  // 2. registro-frontend/package.json
  if (fs.existsSync(files.frontendPkg)) {
    const pkg = JSON.parse(fs.readFileSync(files.frontendPkg, 'utf-8'))
    pkg.version = newVersion
    fs.writeFileSync(files.frontendPkg, JSON.stringify(pkg, null, 2) + '\n', 'utf-8')
    console.log(`  ✓ Updated ${path.relative(rootDir, files.frontendPkg)} -> ${newVersion}`)
  }

  // 3. registro-backend/pkg/version/version.go
  if (fs.existsSync(files.backendVersionGo)) {
    let goContent = fs.readFileSync(files.backendVersionGo, 'utf-8')
    goContent = goContent.replace(/Version\s*=\s*"[^"]+"/, `Version = "${newVersion}"`)
    fs.writeFileSync(files.backendVersionGo, goContent, 'utf-8')
    console.log(`  ✓ Updated ${path.relative(rootDir, files.backendVersionGo)} -> ${newVersion}`)
  }

  // 4. publiccode.yml
  if (fs.existsSync(files.publiccode)) {
    let pc = fs.readFileSync(files.publiccode, 'utf-8')
    pc = pc.replace(/softwareVersion:\s*[^\r\n]+/, `softwareVersion: ${newVersion}`)
    pc = pc.replace(/releaseDate:\s*"[^\r\n]+"/, `releaseDate: "${today}"`)
    fs.writeFileSync(files.publiccode, pc, 'utf-8')
    console.log(`  ✓ Updated ${path.relative(rootDir, files.publiccode)} -> ${newVersion} (releaseDate: ${today})`)
  }

  // 5. CHANGELOG.md
  if (fs.existsSync(files.changelog)) {
    let cl = fs.readFileSync(files.changelog, 'utf-8')
    const headerRegex = new RegExp(`##\\s*\\[${newVersion.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}\\]`)
    if (!headerRegex.test(cl)) {
      const template = `\n## [${newVersion}] — ${today}\n\n### Modifiche\n\n- Rilascio versione ${newVersion}.\n`
      // Insert after top title / description
      const insertPos = cl.indexOf('## [')
      if (insertPos !== -1) {
        cl = cl.slice(0, insertPos) + `## [${newVersion}] — ${today}\n\n### Modifiche\n\n- Note di rilascio per v${newVersion}.\n\n` + cl.slice(insertPos)
      } else {
        cl += template
      }
      fs.writeFileSync(files.changelog, cl, 'utf-8')
      console.log(`  ✓ Added release section to ${path.relative(rootDir, files.changelog)} for [${newVersion}]`)
    }
  }
}

function main() {
  const args = process.argv.slice(2)
  if (args.includes('--check') || args.includes('-c')) {
    checkVersions()
    return
  }

  const bumpArg = args.find(a => !a.startsWith('-')) || 'patch'
  const doGit = args.includes('--git') || args.includes('--tag')

  const currentVersion = readCurrentVersion()
  const nextVersion = computeNextVersion(currentVersion, bumpArg)

  console.log(`🚀 Bumping version: v${currentVersion} ➔ v${nextVersion} (type: ${bumpArg})\n`)
  updateFiles(nextVersion)

  console.log(`\n🎉 Success! All frontend, backend, and metadata files set to v${nextVersion}.\n`)

  if (doGit) {
    try {
      console.log('📦 Creating Git commit & tag...')
      execSync(`git add VERSION registro-frontend/package.json registro-backend/pkg/version/version.go publiccode.yml CHANGELOG.md`, { cwd: rootDir, stdio: 'inherit' })
      execSync(`git commit -m "chore(release): v${nextVersion}"`, { cwd: rootDir, stdio: 'inherit' })
      execSync(`git tag -a "v${nextVersion}" -m "Release v${nextVersion}"`, { cwd: rootDir, stdio: 'inherit' })
      console.log(`🏷️  Created git tag v${nextVersion}`)
    } catch (e) {
      console.error('⚠️  Git commit/tag skipped or failed:', e.message)
    }
  } else {
    console.log('💡 Tip: to commit and create a git tag, run:')
    console.log(`   git commit -am "chore(release): v${nextVersion}"`)
    console.log(`   git tag -a v${nextVersion} -m "Release v${nextVersion}"`)
  }
}

main()
