const API = window.location.origin

const archetypeClasses = {
    'Ассистент': 'badge-assistman',
    'Снайпер': 'badge-sniper',
    'Бомбардир': 'badge-pointer',
    'Защитник-стена': 'badge-iron-defenseman',
    'Атакующий защитник': 'badge-offensive-defenseman',
    'Нарушитель': 'badge-offender',
    'Силовик': 'badge-grinder',
    'Специалист по вбрасываниям': 'badge-faceoff-specialist',
}


const NHL_TEAM_NAMES = {
    ANA: 'Anaheim Ducks',
    ARI: 'Arizona Coyotes',
    BOS: 'Boston Bruins',
    BUF: 'Buffalo Sabres',
    CAR: 'Carolina Hurricanes',
    CBJ: 'Columbus Blue Jackets',
    CGY: 'Calgary Flames',
    CHI: 'Chicago Blackhawks',
    COL: 'Colorado Avalanche',
    DAL: 'Dallas Stars',
    DET: 'Detroit Red Wings',
    EDM: 'Edmonton Oilers',
    FLA: 'Florida Panthers',
    LAK: 'Los Angeles Kings',
    MIN: 'Minnesota Wild',
    MTL: 'Montréal Canadiens',
    NJD: 'New Jersey Devils',
    NSH: 'Nashville Predators',
    NYI: 'New York Islanders',
    NYR: 'New York Rangers',
    OTT: 'Ottawa Senators',
    PHI: 'Philadelphia Flyers',
    PIT: 'Pittsburgh Penguins',
    SEA: 'Seattle Kraken',
    SJS: 'San Jose Sharks',
    STL: 'St. Louis Blues',
    TBL: 'Tampa Bay Lightning',
    TOR: 'Toronto Maple Leafs',
    UTA: 'Utah Mammoth',
    VAN: 'Vancouver Canucks',
    VGK: 'Vegas Golden Knights',
    WPG: 'Winnipeg Jets',
    WSH: 'Washington Capitals',
}

function normalizeTeamCode(team) {
    return String(team || '').trim().toUpperCase()
}

function getTeamFullName(team) {
    const code = normalizeTeamCode(team)
    return NHL_TEAM_NAMES[code] || code || 'Команда не указана'
}

function qs(selector) {
    return document.querySelector(selector)
}

function qsa(selector) {
    return [...document.querySelectorAll(selector)]
}

function escapeHTML(value) {
    return String(value ?? '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#039;')
}

function formatNumber(value, digits = 1) {
    const number = Number(value)
    if (!Number.isFinite(number)) {
        return '—'
    }
    return number.toFixed(digits)
}

function getScoreClass(score) {
    const value = Number(score)
    if (value >= 85) return 'score-elite'
    if (value >= 70) return 'score-good'
    if (value >= 55) return 'score-average'
    if (value >= 40) return 'score-bad'
    return 'score-awful'
}

function getArchetypeClass(archetype) {
    return archetypeClasses[String(archetype ?? '').trim()] || 'badge-default'
}

function seasonStartValue(season) {
    const value = String(season ?? '')
    if (/^\d{8}$/.test(value)) {
        return Number(value.slice(0, 4))
    }

    const [start] = value.split('/')
    const shortYear = Number(start)
    if (!Number.isFinite(shortYear)) {
        return 0
    }

    return shortYear >= 90 ? 1900 + shortYear : 2000 + shortYear
}


function getCurrentNHLSeason() {
    const now = new Date()
    let startYear = now.getFullYear()
    if (now.getMonth() < 8) {
        startYear -= 1
    }
    return `${String(startYear).slice(-2)}/${String(startYear + 1).slice(-2)}`
}

function buildFallbackSeasons() {
    const current = getCurrentNHLSeason()
    const currentStart = seasonStartValue(current)
    const seasons = []
    for (let year = currentStart; year >= 2010; year--) {
        seasons.push(`${String(year).slice(-2)}/${String(year + 1).slice(-2)}`)
    }
    return seasons
}

function sortSeasonsDesc(seasons) {
    return [...new Set(seasons.filter(Boolean))]
        .sort((a, b) => seasonStartValue(b) - seasonStartValue(a))
}

function sortSeasonsAsc(seasons) {
    return [...new Set(seasons.filter(Boolean))]
        .sort((a, b) => seasonStartValue(a) - seasonStartValue(b))
}

function renderState(element, className, message) {
    if (!element) return
    const content = `<div class="state ${className}">${escapeHTML(message)}</div>`

    if (element.tagName === 'TBODY') {
        element.innerHTML = `<tr><td colspan="20">${content}</td></tr>`
        return
    }

    element.innerHTML = content
}

function setLoading(element, message = 'Загрузка...') {
    renderState(element, 'state-loading', message)
}

function setError(element, message = 'Ошибка загрузки данных') {
    renderState(element, 'state-error', message)
}

function setEmpty(element, message = 'Ничего не найдено') {
    renderState(element, 'state-empty', message)
}

async function apiFetch(path, options = {}) {
    const response = await fetch(`${API}${path}`, options)
    const contentType = response.headers.get('content-type') || ''
    const body = contentType.includes('application/json')
        ? await response.json()
        : await response.text()

    if (!response.ok) {
        const message = typeof body === 'object' && body !== null
            ? body.error || body.message || JSON.stringify(body)
            : body

        throw new Error(message || `HTTP ${response.status}`)
    }

    return body
}

async function loadSeasonOptions(select, selected = '') {
    let seasons = []

    try {
        seasons = await apiFetch('/seasons')
    } catch (error) {
        console.warn('Failed to load seasons from backend, using fallback list', error)
    }

    seasons = sortSeasonsDesc([
        ...buildFallbackSeasons(),
        ...seasons,
    ])

    select.innerHTML = ''

    seasons.forEach(season => {
        const option = document.createElement('option')
        option.value = season
        option.textContent = season
        select.appendChild(option)
    })

    if (selected && seasons.includes(selected)) {
        select.value = selected
    } else if (seasons.length > 0) {
        select.value = seasons[0]
    }

    return seasons
}

function updateURLParam(key, value) {
    const url = new URL(window.location.href)
    if (value) {
        url.searchParams.set(key, value)
    } else {
        url.searchParams.delete(key)
    }
    window.history.replaceState({}, '', url)
}

function getTeamLogo(team) {
    const code = normalizeTeamCode(team)
    if (!code) {
        return ''
    }
    return `https://assets.nhle.com/logos/nhl/svg/${encodeURIComponent(code)}_dark.svg`
}

function getPlayerPhoto(nhlId) {
    const id = Number(nhlId)
    if (!Number.isFinite(id) || id <= 0) {
        return ''
    }
    return `https://assets.nhle.com/mugs/nhl/latest/${id}.png`
}

function teamLogoHTML(team, className = 'team-logo') {
    const src = getTeamLogo(team)
    if (!src) {
        return `<span class="${className} team-logo-placeholder">?</span>`
    }
    return `<img class="${className}" src="${src}" alt="${escapeHTML(team)}" loading="lazy" onerror="this.replaceWith(Object.assign(document.createElement('span'), {className: '${className} team-logo-placeholder', textContent: '?' }))">`
}

function playerPhotoHTML(nhlId, playerName, className = 'player-photo') {
    const src = getPlayerPhoto(nhlId)
    const initials = String(playerName || '?')
        .trim()
        .split(/\s+/)
        .slice(0, 2)
        .map(part => part[0] || '')
        .join('')
        .toUpperCase() || '?'

    if (!src) {
        return `<span class="${className} player-photo-placeholder">${escapeHTML(initials)}</span>`
    }

    return `<img class="${className}" src="${src}" alt="${escapeHTML(playerName)}" loading="lazy" onerror="this.replaceWith(Object.assign(document.createElement('span'), {className: '${className} player-photo-placeholder', textContent: '${escapeHTML(initials)}' }))">`
}


window.getTeamLogo = getTeamLogo
window.getPlayerPhoto = getPlayerPhoto
window.teamLogoHTML = teamLogoHTML
window.playerPhotoHTML = playerPhotoHTML
window.getTeamFullName = getTeamFullName
window.normalizeTeamCode = normalizeTeamCode

function downloadCSV(filename, rows) {
    const csv = rows
        .map(row => row
            .map(value => `"${String(value ?? '').replaceAll('"', '""')}"`)
            .join(','))
        .join('\n')

    const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(url)
}
