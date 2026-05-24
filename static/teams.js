const seasonSelect = qs('#seasonSelect')
const teamsGrid = qs('#teamsGrid')
let teams = []

function renderTeams(data) {
    if (!teamsGrid) return
    if (!data.length) {
        setEmpty(teamsGrid, 'Команды не найдены')
        return
    }

    teamsGrid.innerHTML = data.map(team => `
        <a class="team-tile" href="/team/${encodeURIComponent(team.team)}?season=${encodeURIComponent(team.season)}">
            <div class="team-card-header compact">
                ${teamLogoHTML(team.team, 'team-logo-medium')}
                <div>
                    <h3>${escapeHTML(getTeamFullName(team.team))}</h3>
                    <div class="subtitle"><span class="team-code-pill">${escapeHTML(team.team)}</span> ${team.players_count} игроков</div>
                </div>
            </div>
            <div class="stat-grid mini-grid">
                <div><div class="stat-title">Avg Overall</div><div class="stat-value ${getScoreClass(team.average_overall)}">${formatNumber(team.average_overall)}</div></div>
                <div><div class="stat-title">Best</div><div class="stat-value small-stat">${escapeHTML(team.best_player || '—')}</div></div>
                <div><div class="stat-title">Очки</div><div class="stat-value">${team.points}</div></div>
                <div><div class="stat-title">Архетип</div><div><span class="badge ${getArchetypeClass(team.top_archetype)}">${escapeHTML(team.top_archetype || '—')}</span></div></div>
            </div>
        </a>
    `).join('')
}

async function loadTeams() {
    const season = seasonSelect.value
    updateURLParam('season', season)
    setLoading(teamsGrid, 'Загружаю команды...')
    try {
        teams = await apiFetch(`/api/teams?season=${encodeURIComponent(season)}`)
        renderTeams(teams)
    } catch (error) {
        setError(teamsGrid, error.message)
    }
}

function exportTeamsCsv() {
    const rows = [
        ['Team', 'FullName', 'Players', 'AverageOverall', 'AverageBetterThanPercent', 'BestPlayer', 'BestOverall', 'Goals', 'Assists', 'Points', 'TopArchetype'],
        ...teams.map(t => [t.team, getTeamFullName(t.team), t.players_count, t.average_overall, t.average_percentile, t.best_player, t.best_overall, t.goals, t.assists, t.points, t.top_archetype]),
    ]
    downloadCSV(`teams_${seasonSelect.value.replace('/', '-')}.csv`, rows)
}

async function init() {
    const selected = new URLSearchParams(window.location.search).get('season') || ''
    await loadSeasonOptions(seasonSelect, selected)
    seasonSelect.addEventListener('change', loadTeams)
    qs('#exportTeamsCsv')?.addEventListener('click', exportTeamsCsv)
    await loadTeams()
}

init().catch(error => setError(teamsGrid, error.message))
