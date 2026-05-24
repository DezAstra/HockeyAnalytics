let allPlayers = []
let currentSort = 'overall_score'
let sortDirection = 'desc'
let lastFilteredPlayers = []

const table = qs('#playersTable')
const seasonSelect = qs('#seasonSelect')
const searchInput = qs('#searchInput')
const summaryCards = qs('#summaryCards')
const resultCount = qs('#resultCount')

const POSITION_GROUPS = {
    C: ['C'],
    L: ['L', 'LW'],
    LW: ['L', 'LW'],
    R: ['R', 'RW'],
    RW: ['R', 'RW'],
    D: ['D'],
}

function isAllValue(value) {
    return !value || value === 'all'
}

function normalizePosition(position) {
    return String(position || '').trim().toUpperCase()
}

function matchesPosition(playerPosition, selectedPosition) {
    if (isAllValue(selectedPosition)) return true
    const normalizedSelectedPosition = normalizePosition(selectedPosition)
    const allowedPositions = POSITION_GROUPS[normalizedSelectedPosition] || [normalizedSelectedPosition]
    return allowedPositions.includes(normalizePosition(playerPosition))
}

function renderSummary(players) {
    if (!summaryCards) return
    if (!players.length) {
        summaryCards.innerHTML = ''
        return
    }

    const avgNormalized = players.reduce((sum, p) => sum + Number(p.normalized_score || 0), 0) / players.length
    const avgContext = players.reduce((sum, p) => sum + Number(p.context_score || 0), 0) / players.length
    const maxOverall = Math.max(...players.map(p => Number(p.overall_score || 0)))
    const avgPercentile = players.reduce((sum, p) => sum + Number(p.overall_percentile || 0), 0) / players.length
    const topPlayer = [...players].sort((a, b) => Number(b.overall_score || 0) - Number(a.overall_score || 0))[0]

    summaryCards.innerHTML = `
        <div class="stat-box"><div class="stat-title">Игроков</div><div class="stat-value">${players.length}</div></div>
        <div class="stat-box"><div class="stat-title">Лучший Overall</div><div class="stat-value ${getScoreClass(maxOverall)}">${formatNumber(maxOverall)}</div></div>
        <div class="stat-box"><div class="stat-title">Средний Лучше %</div><div class="stat-value">${formatNumber(avgPercentile)}</div></div>
        <div class="stat-box"><div class="stat-title">Средний Normalized</div><div class="stat-value ${getScoreClass(avgNormalized)}">${formatNumber(avgNormalized)}</div></div>
        <div class="stat-box"><div class="stat-title">Средний Context</div><div class="stat-value ${getScoreClass(avgContext)}">${formatNumber(avgContext)}</div></div>
        <div class="stat-box">
            <div class="stat-title">Лучший игрок</div>
            <div class="stat-value small-stat">
                <a class="player-cell" href="/player?id=${encodeURIComponent(topPlayer.player_id)}">
                    ${playerPhotoHTML(topPlayer.nhl_id, topPlayer.player, 'player-photo-small')}
                    <span>${escapeHTML(topPlayer.player)}</span>
                </a>
            </div>
        </div>
    `
}

function renderPlayers(players) {
    if (!table) return
    table.innerHTML = ''
    if (resultCount) resultCount.textContent = `${players.length} игроков`

    if (!players.length) {
        table.innerHTML = `<tr><td colspan="10"><div class="state state-empty">По текущим фильтрам игроков нет</div></td></tr>`
        return
    }

    table.innerHTML = players.map(player => `
        <tr>
            <td>
                <a class="player-cell" href="/player?id=${encodeURIComponent(player.player_id)}">
                    ${playerPhotoHTML(player.nhl_id, player.player, 'player-photo-small')}
                    <span>${escapeHTML(player.player)}</span>
                </a>
            </td>
            <td>
                <a class="team-cell" href="/team/${encodeURIComponent(player.team)}?season=${encodeURIComponent(player.season || seasonSelect.value)}">
                    ${teamLogoHTML(player.team, 'team-logo-small')}
                    <span>${escapeHTML(player.team || '—')}</span>
                </a>
            </td>
            <td>${escapeHTML(player.position || '—')}</td>
            <td>${player.games_played ?? 0}</td>
            <td>${player.points ?? 0}</td>
            <td class="${getScoreClass(player.normalized_score)}">${formatNumber(player.normalized_score)}</td>
            <td class="${getScoreClass(player.context_score)}">${formatNumber(player.context_score)}</td>
            <td class="${getScoreClass(player.overall_score)}">${formatNumber(player.overall_score)}</td>
            <td>${formatNumber(player.overall_percentile)}</td>
            <td><span class="badge ${getArchetypeClass(player.archetype)}">${escapeHTML(player.archetype || 'Баланс')}</span></td>
        </tr>
    `).join('')
}

function populateTeamFilter(data) {
    const select = qs('#teamFilter')
    if (!select) return
    const current = select.value
    const teams = [...new Set(data.map(p => p.team).filter(Boolean))].sort()
    select.innerHTML = '<option value="">Все команды</option>'
    teams.forEach(team => {
        const option = document.createElement('option')
        option.value = team
        option.textContent = `${team} — ${getTeamFullName(team)}`
        select.appendChild(option)
    })
    if (teams.includes(current)) select.value = current
}

function getFilteredPlayers() {
    const search = searchInput ? searchInput.value.trim().toLowerCase() : ''
    const position = qs('#positionFilter')?.value || ''
    const team = qs('#teamFilter')?.value || ''
    const archetype = qs('#archetypeFilter')?.value || ''
    const minGames = Number(qs('#minGamesFilter')?.value || 0)

    let filtered = [...allPlayers]

    if (search) {
        filtered = filtered.filter(player => {
            const name = String(player.player || '').toLowerCase()
            const teamCode = String(player.team || '').toLowerCase()
            const teamName = getTeamFullName(player.team).toLowerCase()
            return name.includes(search) || teamCode.includes(search) || teamName.includes(search)
        })
    }

    if (!isAllValue(position)) filtered = filtered.filter(player => matchesPosition(player.position, position))
    if (!isAllValue(team)) filtered = filtered.filter(player => player.team === team)
    if (!isAllValue(archetype)) filtered = filtered.filter(player => player.archetype === archetype)
    if (minGames > 0) filtered = filtered.filter(player => Number(player.games_played || 0) >= minGames)

    if (currentSort) {
        filtered.sort((a, b) => {
            const aVal = a[currentSort]
            const bVal = b[currentSort]
            if (typeof aVal === 'string' || typeof bVal === 'string') {
                return sortDirection === 'asc'
                    ? String(aVal || '').localeCompare(String(bVal || ''), 'ru')
                    : String(bVal || '').localeCompare(String(aVal || ''), 'ru')
            }
            return sortDirection === 'asc'
                ? Number(aVal || 0) - Number(bVal || 0)
                : Number(bVal || 0) - Number(aVal || 0)
        })
    }

    return filtered
}

function applyFilters() {
    lastFilteredPlayers = getFilteredPlayers()
    renderSummary(lastFilteredPlayers)
    renderPlayers(lastFilteredPlayers)
}

async function loadAnalytics() {
    const season = seasonSelect?.value
    if (!season) return
    updateURLParam('season', season)
    setLoading(table, 'Загружаю сезон. Если его нет в базе, backend автоматически импортирует его из NHL API...')

    try {
        allPlayers = await apiFetch(`/analytics?season=${encodeURIComponent(season)}`)
        populateTeamFilter(allPlayers)
        applyFilters()
    } catch (error) {
        setError(table, error.message)
        if (summaryCards) summaryCards.innerHTML = ''
        if (resultCount) resultCount.textContent = ''
    }
}

function exportCurrentTable() {
    const rows = [
        ['Player', 'Team', 'Position', 'GP', 'Goals', 'Assists', 'Points', 'Normalized', 'Context', 'Overall', 'BetterThanPercent', 'Archetype'],
        ...lastFilteredPlayers.map(p => [
            p.player,
            p.team,
            p.position,
            p.games_played,
            p.goals,
            p.assists,
            p.points,
            p.normalized_score,
            p.context_score,
            p.overall_score,
            p.overall_percentile,
            p.archetype,
        ]),
    ]
    downloadCSV(`leaderboard_${seasonSelect.value.replace('/', '-')}.csv`, rows)
}

function bindEvents() {
    if (seasonSelect) seasonSelect.addEventListener('change', loadAnalytics)
    if (searchInput) searchInput.addEventListener('input', applyFilters)

    ;['#positionFilter', '#teamFilter', '#archetypeFilter', '#minGamesFilter'].forEach(selector => {
        const element = qs(selector)
        if (element) element.addEventListener('change', applyFilters)
    })

    qs('#resetFilters')?.addEventListener('click', () => {
        if (searchInput) searchInput.value = ''
        if (qs('#positionFilter')) qs('#positionFilter').value = 'all'
        if (qs('#teamFilter')) qs('#teamFilter').value = ''
        if (qs('#archetypeFilter')) qs('#archetypeFilter').value = ''
        if (qs('#minGamesFilter')) qs('#minGamesFilter').value = '0'
        applyFilters()
    })

    qs('#exportCsv')?.addEventListener('click', exportCurrentTable)

    qsa('th[data-sort]').forEach(th => {
        th.addEventListener('click', () => {
            const field = th.dataset.sort
            if (currentSort === field) {
                sortDirection = sortDirection === 'asc' ? 'desc' : 'asc'
            } else {
                currentSort = field
                sortDirection = 'desc'
            }
            applyFilters()
        })
    })
}

async function init() {
    bindEvents()
    const selected = new URLSearchParams(window.location.search).get('season') || ''
    try {
        await loadSeasonOptions(seasonSelect, selected)
        await loadAnalytics()
    } catch (error) {
        setError(table, error.message)
        if (summaryCards) summaryCards.innerHTML = ''
        if (resultCount) resultCount.textContent = ''
    }
}

init()
