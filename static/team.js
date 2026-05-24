const pathParts = window.location.pathname.split('/')
const team = decodeURIComponent(pathParts[pathParts.length - 1] || '')

const table = qs('#playersTable')
const seasonSelect = qs('#seasonSelect')
const teamHeader = qs('#teamHeader')
const teamSummary = qs('#teamSummary')
const teamProfile = qs('#teamProfile')
const teamHistorySummary = qs('#teamHistorySummary')
const teamHistoryBreakdown = qs('#teamHistoryBreakdown')
const teamHistoryChartCanvas = qs('#teamHistoryChart')

let currentSort = 'overall_score'
let sortDirection = 'desc'
let currentPlayers = []
let currentTeamData = null
let currentTeamHistory = null
let teamHistoryChart = null

function parseSeason(season) {
    const value =
        String(season || '').trim()

    if (!value) {
        return 0
    }

    if (value.includes('/')) {
        const parts =
            value.split('/')

        const start =
            Number(parts[0])

        if (Number.isNaN(start)) {
            return 0
        }

        return start >= 70
            ? 1900 + start
            : 2000 + start
    }

    if (value.length === 8) {
        return Number(value.slice(0, 4)) || 0
    }

    return Number(value) || 0
}

function getSortedPlayers(players) {
    const sorted = [...players]

    if (!currentSort) {
        return sorted
    }

    sorted.sort((a, b) => {
        const aValue = a[currentSort]
        const bValue = b[currentSort]

        if (typeof aValue === 'string' || typeof bValue === 'string') {
            return sortDirection === 'asc'
                ? String(aValue || '').localeCompare(String(bValue || ''), 'ru')
                : String(bValue || '').localeCompare(String(aValue || ''), 'ru')
        }

        return sortDirection === 'asc'
            ? Number(aValue || 0) - Number(bValue || 0)
            : Number(bValue || 0) - Number(aValue || 0)
    })

    return sorted
}

function signedNumber(value) {
    const number = Number(value)

    if (!Number.isFinite(number)) {
        return '—'
    }

    if (number > 0) {
        return `+${formatNumber(number)}`
    }

    return formatNumber(number)
}

function sortHistoryAsc(history) {
    return [...history].sort((a, b) =>
        parseSeason(a.season) - parseSeason(b.season),
    )
}

function calculateThreeSeasonsBackTrend(history) {
    const sorted = sortHistoryAsc(history)

    if (sorted.length < 4) {
        return {
            trend: 'unknown',
            delta: null,
            currentSeason: sorted.at(-1)?.season || '',
            compareSeason: '',
            description: 'Нужно минимум 4 сезона для сравнения с сезоном 3 года назад',
        }
    }

    const current = sorted[sorted.length - 1]
    const compare = sorted[sorted.length - 4]

    const delta =
        Number(current.average_overall || 0) -
        Number(compare.average_overall || 0)

    let trend = 'stable'

    if (delta > 1.5) {
        trend = 'up'
    } else if (delta < -1.5) {
        trend = 'down'
    }

    return {
        trend,
        delta,
        currentSeason: current.season,
        compareSeason: compare.season,
        description: `${signedNumber(delta)} Overall к сезону ${compare.season}`,
    }
}

function getTrendLabel(trend) {
    if (trend === 'up') {
        return 'Команда усиливается'
    }

    if (trend === 'down') {
        return 'Команда проседает'
    }

    if (trend === 'stable') {
        return 'Команда стабильна'
    }

    return 'Недостаточно сезонов'
}

function getTrendClass(trend) {
    if (trend === 'up') {
        return 'trend-up-text'
    }

    if (trend === 'down') {
        return 'trend-down-text'
    }

    return ''
}

function renderTeamHeader(data) {
    if (!teamHeader) {
        return
    }

    const teamCode = normalizeTeamCode(data.team)
    const teamName = getTeamFullName(data.team)

    teamHeader.innerHTML = `
        <div class="team-card-header">
            ${teamLogoHTML(teamCode, 'team-logo-large')}
            <div>
                <div class="title">${escapeHTML(teamName)}</div>
                <div class="subtitle">
                    <span class="team-code-pill">${escapeHTML(teamCode || '—')}</span>
                    Сезон ${escapeHTML(data.season)}
                </div>
            </div>
        </div>
    `
}

function getMostCommon(values) {
    const counts = {}

    values
        .filter(Boolean)
        .forEach(value => {
            counts[value] = (counts[value] || 0) + 1
        })

    return Object
        .entries(counts)
        .sort((a, b) => b[1] - a[1])[0]?.[0] || '—'
}

function renderSummary(players) {
    if (!teamSummary) {
        return
    }

    if (!players.length) {
        teamSummary.innerHTML = ''
        return
    }

    const goals = players.reduce((sum, p) => sum + Number(p.goals || 0), 0)
    const assists = players.reduce((sum, p) => sum + Number(p.assists || 0), 0)
    const points = players.reduce((sum, p) => sum + Number(p.points || 0), 0)
    const avgOverall = players.reduce((sum, p) => sum + Number(p.overall_score || 0), 0) / players.length
    const topPlayer = [...players].sort((a, b) => Number(b.overall_score || 0) - Number(a.overall_score || 0))[0]
    const topGoal = [...players].sort((a, b) => Number(b.goals || 0) - Number(a.goals || 0))[0]
    const topAssist = [...players].sort((a, b) => Number(b.assists || 0) - Number(a.assists || 0))[0]

    teamSummary.innerHTML = `
        <div class="stat-box"><div class="stat-title">Игроков</div><div class="stat-value">${players.length}</div></div>
        <div class="stat-box"><div class="stat-title">Голы</div><div class="stat-value">${goals}</div></div>
        <div class="stat-box"><div class="stat-title">Ассисты</div><div class="stat-value">${assists}</div></div>
        <div class="stat-box"><div class="stat-title">Очки</div><div class="stat-value">${points}</div></div>
        <div class="stat-box"><div class="stat-title">Средний Overall</div><div class="stat-value ${getScoreClass(avgOverall)}">${formatNumber(avgOverall)}</div></div>
        <div class="stat-box"><div class="stat-title">Лучший игрок</div><div class="stat-value small-stat"><a href="/player?id=${topPlayer.player_id}">${escapeHTML(topPlayer.player)}</a></div></div>
        <div class="stat-box"><div class="stat-title">Лучший снайпер</div><div class="stat-value small-stat">${escapeHTML(topGoal.player)} · ${topGoal.goals}</div></div>
        <div class="stat-box"><div class="stat-title">Лучший ассистент</div><div class="stat-value small-stat">${escapeHTML(topAssist.player)} · ${topAssist.assists}</div></div>
    `
}

function renderTeamProfile(players) {
    if (!teamProfile || !players.length) {
        return
    }

    const positions = getMostCommon(players.map(p => p.position))
    const archetype = getMostCommon(players.map(p => p.archetype))
    const offensePlayers = players.filter(p => ['Снайпер', 'Ассистент', 'Бомбардир', 'Атакующий защитник'].includes(p.archetype)).length
    const physicalPlayers = players.filter(p => ['Силовик', 'Защитник-стена', 'Нарушитель'].includes(p.archetype)).length
    const profileText = offensePlayers >= physicalPlayers
        ? 'Профиль команды смещён в сторону атакующих и созидательных ролей.'
        : 'Профиль команды смещён в сторону силовой и оборонительной игры.'

    teamProfile.innerHTML = `
        <div class="explain-panel">
            <h3>Профиль команды</h3>
            <div class="grid-2">
                <div>
                    <div class="stat-title">Главный архетип</div>
                    <div><span class="badge ${getArchetypeClass(archetype)}">${escapeHTML(archetype)}</span></div>
                </div>
                <div>
                    <div class="stat-title">Самая частая позиция</div>
                    <div class="subtitle">${escapeHTML(positions)}</div>
                </div>
            </div>
            <p class="subtitle">${escapeHTML(profileText)}</p>
        </div>
    `
}

function renderTeamHistorySummary(data) {
    if (!teamHistorySummary) {
        return
    }

    const history = sortHistoryAsc(data?.history || [])

    if (!history.length) {
        setEmpty(teamHistorySummary, 'Истории команды пока нет')
        return
    }

    const last = history[history.length - 1]
    const best = [...history].sort((a, b) => Number(b.average_overall || 0) - Number(a.average_overall || 0))[0]
    const trendInfo = calculateThreeSeasonsBackTrend(history)
    const trendLabel = getTrendLabel(trendInfo.trend)
    const trendClass = getTrendClass(trendInfo.trend)

    teamHistorySummary.innerHTML = `
        <div class="stat-box">
            <div class="stat-title">Сезонов в базе</div>
            <div class="stat-value">${history.length}</div>
        </div>

        <div class="stat-box">
            <div class="stat-title">Тренд состава</div>
            <div class="stat-value small-stat ${trendClass}">
                ${escapeHTML(trendLabel)}
            </div>
            <div class="subtitle">${escapeHTML(trendInfo.description)}</div>
        </div>

        <div class="stat-box">
            <div class="stat-title">Лучший сезон</div>
            <div class="stat-value ${getScoreClass(best.average_overall)}">
                ${formatNumber(best.average_overall)}
            </div>
            <div class="subtitle">${escapeHTML(best.season)}</div>
        </div>

        <div class="stat-box">
            <div class="stat-title">Лучший игрок последнего сезона</div>
            <div class="stat-value small-stat">
                <a href="/player?id=${encodeURIComponent(last.best_player_id)}">
                    ${escapeHTML(last.best_player || '—')}
                </a>
            </div>
            <div class="subtitle">Overall ${formatNumber(last.best_overall)}</div>
        </div>
    `
}

function renderTeamHistoryChart(history) {
    if (!teamHistoryChartCanvas || typeof Chart === 'undefined') {
        return
    }

    if (teamHistoryChart) {
        teamHistoryChart.destroy()
    }

    if (!history.length) {
        return
    }

    const sortedHistory = sortHistoryAsc(history)

    teamHistoryChart = new Chart(teamHistoryChartCanvas, {
        type: 'line',
        data: {
            labels: sortedHistory.map(item => item.season),
            datasets: [
                {
                    label: 'Avg Overall',
                    data: sortedHistory.map(item => Number(item.average_overall || 0)),
                    tension: 0.3,
                },
                {
                    label: 'Avg Context',
                    data: sortedHistory.map(item => Number(item.average_context || 0)),
                    tension: 0.3,
                },
                {
                    label: 'Avg Normalized',
                    data: sortedHistory.map(item => Number(item.average_normalized || 0)),
                    tension: 0.3,
                },
            ],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            interaction: {
                mode: 'index',
                intersect: false,
            },
            plugins: {
                legend: {
                    labels: {
                        color: '#edf5ff',
                    },
                },
                tooltip: {
                    callbacks: {
                        label(context) {
                            return `${context.dataset.label}: ${formatNumber(context.raw)}`
                        },
                    },
                },
            },
            scales: {
                x: {
                    ticks: {
                        color: '#95a4bb',
                    },
                    grid: {
                        color: 'rgba(149,164,187,.12)',
                    },
                },
                y: {
                    ticks: {
                        color: '#95a4bb',
                    },
                    grid: {
                        color: 'rgba(149,164,187,.12)',
                    },
                },
            },
        },
    })
}

function renderTeamHistoryBreakdown(history) {
    if (!teamHistoryBreakdown) {
        return
    }

    if (!history.length) {
        teamHistoryBreakdown.innerHTML = ''
        return
    }

    const sortedHistory = sortHistoryAsc(history)

    teamHistoryBreakdown.innerHTML = `
        <div class="table-wrapper history-table-wrapper">
            <table>
                <thead>
                    <tr>
                        <th>Сезон</th>
                        <th>Игроков</th>
                        <th>Avg Overall</th>
                        <th>Avg Context</th>
                        <th>Avg Normalized</th>
                        <th>Лучший игрок</th>
                        <th>Очки</th>
                    </tr>
                </thead>
                <tbody>
                    ${sortedHistory.map(item => `
                        <tr>
                            <td>${escapeHTML(item.season)}</td>
                            <td>${item.players_count}</td>
                            <td class="${getScoreClass(item.average_overall)}">${formatNumber(item.average_overall)}</td>
                            <td class="${getScoreClass(item.average_context)}">${formatNumber(item.average_context)}</td>
                            <td class="${getScoreClass(item.average_normalized)}">${formatNumber(item.average_normalized)}</td>
                            <td>
                                <a class="player-cell" href="/player?id=${encodeURIComponent(item.best_player_id)}">
                                    ${playerPhotoHTML(item.best_player_nhl_id, item.best_player, 'player-photo-tiny')}
                                    <span>${escapeHTML(item.best_player || '—')}</span>
                                </a>
                            </td>
                            <td>${item.points}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `
}

function renderTeamHistory(data) {
    currentTeamHistory = data
    const history = sortHistoryAsc(data?.history || [])

    renderTeamHistorySummary({
        ...data,
        history,
    })
    renderTeamHistoryChart(history)
    renderTeamHistoryBreakdown(history)
}

function renderPlayers(players) {
    if (!table) {
        return
    }

    const sortedPlayers = getSortedPlayers(players)

    if (!sortedPlayers.length) {
        setEmpty(table, 'По выбранному сезону игроков нет')
        return
    }

    table.innerHTML = sortedPlayers.map(player => `
        <tr>
            <td>
                <a class="player-cell" href="/player?id=${encodeURIComponent(player.player_id)}">
                    ${playerPhotoHTML(player.nhl_id, player.player, 'player-photo-small')}
                    <span>${escapeHTML(player.player)}</span>
                </a>
            </td>
            <td>${escapeHTML(player.position || '—')}</td>
            <td>${player.goals ?? 0}</td>
            <td>${player.assists ?? 0}</td>
            <td>${player.points ?? 0}</td>
            <td class="${getScoreClass(player.normalized_score)}">${formatNumber(player.normalized_score)}</td>
            <td class="${getScoreClass(player.context_score)}">${formatNumber(player.context_score)}</td>
            <td class="${getScoreClass(player.overall_score)}">${formatNumber(player.overall_score)}</td>
            <td><span class="badge ${getArchetypeClass(player.archetype)}">${escapeHTML(player.archetype || 'Баланс')}</span></td>
        </tr>
    `).join('')
}

async function loadTeamHistory() {
    if (teamHistorySummary) {
        setLoading(teamHistorySummary, 'Загружаю историю команды...')
    }

    try {
        const data = await apiFetch(`/api/team/${encodeURIComponent(team)}/history`)
        renderTeamHistory(data)
    } catch (error) {
        if (teamHistorySummary) {
            setError(teamHistorySummary, error.message)
        }

        if (teamHistoryBreakdown) {
            teamHistoryBreakdown.innerHTML = ''
        }
    }
}

async function loadTeam() {
    const season = seasonSelect?.value

    if (!season) {
        return
    }

    updateURLParam('season', season)
    setLoading(table, 'Загружаю команду...')

    try {
        const data = await apiFetch(`/api/team/${encodeURIComponent(team)}?season=${encodeURIComponent(season)}`)
        currentTeamData = data
        currentPlayers = data.players || []
        renderTeamHeader(data)
        renderSummary(currentPlayers)
        renderTeamProfile(currentPlayers)
        renderPlayers(currentPlayers)
    } catch (error) {
        setError(table, error.message)

        if (teamSummary) {
            teamSummary.innerHTML = ''
        }

        if (teamProfile) {
            teamProfile.innerHTML = ''
        }
    }
}

function exportTeamCsv() {
    const rows = [
        ['Player', 'Position', 'Goals', 'Assists', 'Points', 'Normalized', 'Context', 'Overall', 'Archetype'],
        ...getSortedPlayers(currentPlayers).map(p => [
            p.player,
            p.position,
            p.goals,
            p.assists,
            p.points,
            p.normalized_score,
            p.context_score,
            p.overall_score,
            p.archetype,
        ]),
    ]

    downloadCSV(
        `team_${team}_${seasonSelect.value.replace('/', '-')}.csv`,
        rows,
    )
}

function exportTeamHistoryCsv() {
    const history = sortHistoryAsc(currentTeamHistory?.history || [])

    const rows = [
        ['Season', 'Players', 'AverageOverall', 'AverageContext', 'AverageNormalized', 'BestPlayer', 'BestOverall', 'Goals', 'Assists', 'Points'],
        ...history.map(item => [
            item.season,
            item.players_count,
            item.average_overall,
            item.average_context,
            item.average_normalized,
            item.best_player,
            item.best_overall,
            item.goals,
            item.assists,
            item.points,
        ]),
    ]

    downloadCSV(
        `team_history_${team}.csv`,
        rows,
    )
}

function bindEvents() {
    if (seasonSelect) {
        seasonSelect.addEventListener('change', loadTeam)
    }

    qs('#exportTeamCsv')?.addEventListener('click', exportTeamCsv)
    qs('#exportTeamHistoryCsv')?.addEventListener('click', exportTeamHistoryCsv)

    qsa('th[data-sort]').forEach(th => {
        th.addEventListener('click', () => {
            const field = th.dataset.sort

            if (currentSort === field) {
                sortDirection = sortDirection === 'asc' ? 'desc' : 'asc'
            } else {
                currentSort = field
                sortDirection = 'desc'
            }

            renderPlayers(currentPlayers)
        })
    })
}

async function init() {
    bindEvents()

    const selected = new URLSearchParams(window.location.search).get('season') || ''

    try {
        await loadSeasonOptions(seasonSelect, selected)
        await Promise.all([
            loadTeam(),
            loadTeamHistory(),
        ])
    } catch (error) {
        setError(table, error.message)
    }
}

init()