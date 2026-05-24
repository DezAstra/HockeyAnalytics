let radarChart = null
let seasonPlayers = []
let selectedPlayer1 = null
let selectedPlayer2 = null

const params = new URLSearchParams(window.location.search)
const seasonSelect = qs('#seasonSelect')
const comparisonCard = qs('#comparisonCard')

function getSelectedSeason() {
    return seasonSelect.value
}

async function loadPlayersForSeason() {
    const season = getSelectedSeason()
    seasonPlayers = await apiFetch(`/analytics?season=${encodeURIComponent(season)}`)
}

function renderSearchResults(target, input, resultsBox, setter) {
    const query = input.value.trim().toLowerCase()

    if (query.length < 2) {
        resultsBox.innerHTML = ''
        return
    }

    const players = seasonPlayers
        .filter(player => player.player.toLowerCase().includes(query))
        .slice(0, 8)

    resultsBox.innerHTML = players.map(player => `
        <div class="compare-result" data-id="${player.player_id}" data-name="${escapeHTML(player.player)}">
            ${playerPhotoHTML(player.nhl_id, player.player, 'player-photo-tiny')}
            <span>${escapeHTML(player.player)} (${escapeHTML(player.team || '—')})</span>
        </div>
    `).join('')

    qsa(`#${target} .compare-result`).forEach(row => {
        row.addEventListener('click', () => {
            setter(row.dataset.id)
            input.value = row.dataset.name
            resultsBox.innerHTML = ''
        })
    })
}

function winnerClass(a, b) {
    return Number(a) > Number(b) ? 'winner' : ''
}

function signedDelta(a, b, digits = 1) {
    const delta = Number(a || 0) - Number(b || 0)
    const formatted = formatNumber(Math.abs(delta), digits)
    if (delta > 0) return `+${formatted}`
    if (delta < 0) return `-${formatted}`
    return '0.0'
}

function buildVerdict(data) {
    const p1 = data.player1
    const p2 = data.player2
    const overallDelta = Number(p1.overall_score || 0) - Number(p2.overall_score || 0)
    const pointsDelta = Number(p1.points || 0) - Number(p2.points || 0)
    const goalsDelta = Number(p1.goals || 0) - Number(p2.goals || 0)

    const leader = overallDelta >= 0 ? p1 : p2
    const opponent = overallDelta >= 0 ? p2 : p1

    const reasons = []
    if (Math.abs(overallDelta) >= 1) reasons.push(`${escapeHTML(leader.player)} выше по overall на ${signedDelta(leader.overall_score, opponent.overall_score)}.`)
    else reasons.push('По overall игроки очень близки.')

    if (Math.abs(pointsDelta) > 0) {
        const pointsLeader = pointsDelta >= 0 ? p1 : p2
        const pointsOpponent = pointsDelta >= 0 ? p2 : p1
        reasons.push(`${escapeHTML(pointsLeader.player)} сильнее по очкам: ${signedDelta(pointsLeader.points, pointsOpponent.points, 0)}.`)
    }

    if (Math.abs(goalsDelta) > 0) {
        const goalsLeader = goalsDelta >= 0 ? p1 : p2
        const goalsOpponent = goalsDelta >= 0 ? p2 : p1
        reasons.push(`${escapeHTML(goalsLeader.player)} сильнее как снайпер: ${signedDelta(goalsLeader.goals, goalsOpponent.goals, 0)} голов.`)
    }

    return reasons.join(' ')
}

function renderComparison(data) {
    comparisonCard.innerHTML = `
        <div class="compare-grid">
            <div class="compare-player">
                ${playerPhotoHTML(data.player1.nhl_id, data.player1.player, 'player-photo-compare')}
                <h2>${escapeHTML(data.player1.player)}</h2>
                <p class="subtitle team-cell inline-team">
                    ${teamLogoHTML(data.player1.team, 'team-logo-small')}
                    <span>${escapeHTML(data.player1.team || '—')} · ${escapeHTML(data.player1.position || '—')}</span>
                </p>
                <div class="stat-value ${getScoreClass(data.player1.overall_score)}">${formatNumber(data.player1.overall_score)}</div>
            </div>
            <div class="compare-player">
                ${playerPhotoHTML(data.player2.nhl_id, data.player2.player, 'player-photo-compare')}
                <h2>${escapeHTML(data.player2.player)}</h2>
                <p class="subtitle team-cell inline-team">
                    ${teamLogoHTML(data.player2.team, 'team-logo-small')}
                    <span>${escapeHTML(data.player2.team || '—')} · ${escapeHTML(data.player2.position || '—')}</span>
                </p>
                <div class="stat-value ${getScoreClass(data.player2.overall_score)}">${formatNumber(data.player2.overall_score)}</div>
            </div>
        </div>

        <div class="explain-panel">
            <h3>Итог сравнения</h3>
            <p class="subtitle">${buildVerdict(data)}</p>
        </div>

        <div class="delta-grid">
            <div class="stat-box"><div class="stat-title">Разница Overall</div><div class="stat-value">${signedDelta(data.player1.overall_score, data.player2.overall_score)}</div></div>
            <div class="stat-box"><div class="stat-title">Разница очков</div><div class="stat-value">${signedDelta(data.player1.points, data.player2.points, 0)}</div></div>
            <div class="stat-box"><div class="stat-title">Разница голов</div><div class="stat-value">${signedDelta(data.player1.goals, data.player2.goals, 0)}</div></div>
        </div>

        <div class="table-wrapper" style="max-height:none;margin-top:18px">
            <table class="compare-table">
                <thead>
                    <tr><th>Показатель</th><th>${escapeHTML(data.player1.player)}</th><th>${escapeHTML(data.player2.player)}</th></tr>
                </thead>
                <tbody>
                    ${comparisonRow('Голы', data.player1.goals, data.player2.goals)}
                    ${comparisonRow('Ассисты', data.player1.assists, data.player2.assists)}
                    ${comparisonRow('Очки', data.player1.points, data.player2.points)}
                    ${comparisonRow('Хиты', data.player1.hits, data.player2.hits)}
                    ${comparisonRow('Блоки', data.player1.blocks, data.player2.blocks)}
                    ${comparisonRow('Base', data.player1.base_score, data.player2.base_score, true)}
                    ${comparisonRow('Normalized', data.player1.normalized_score, data.player2.normalized_score, true)}
                    ${comparisonRow('Context', data.player1.context_score, data.player2.context_score, true)}
                    ${comparisonRow('Overall', data.player1.overall_score, data.player2.overall_score, true)}
                </tbody>
            </table>
        </div>
    `
}

function comparisonRow(label, first, second, fixed = false) {
    return `
        <tr>
            <td>${escapeHTML(label)}</td>
            <td class="${winnerClass(first, second)}">${fixed ? formatNumber(first) : escapeHTML(first)}</td>
            <td class="${winnerClass(second, first)}">${fixed ? formatNumber(second) : escapeHTML(second)}</td>
        </tr>
    `
}

function normalize(value, max) {
    return max === 0 ? 0 : (Number(value || 0) / max) * 100
}

function renderRadar(data) {
    const ctx = qs('#comparisonChart')
    const metrics = [
        ['Голы', 'goals'],
        ['Ассисты', 'assists'],
        ['Хиты', 'hits'],
        ['Блоки', 'blocks'],
        ['Overall', 'overall_score'],
    ]

    const maxValues = Object.fromEntries(metrics.map(([, key]) => [
        key,
        Math.max(Number(data.player1[key] || 0), Number(data.player2[key] || 0)),
    ]))

    if (radarChart) {
        radarChart.destroy()
    }

    radarChart = new Chart(ctx, {
        type: 'radar',
        data: {
            labels: metrics.map(([label]) => label),
            datasets: [
                {
                    label: data.player1.player,
                    data: metrics.map(([, key]) => normalize(data.player1[key], maxValues[key])),
                },
                {
                    label: data.player2.player,
                    data: metrics.map(([, key]) => normalize(data.player2[key], maxValues[key])),
                },
            ],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: { legend: { labels: { color: '#edf5ff' } } },
            scales: {
                r: {
                    min: 0,
                    max: 100,
                    ticks: { display: false },
                    pointLabels: { color: '#edf5ff' },
                    grid: { color: 'rgba(149,164,187,.18)' },
                    angleLines: { color: 'rgba(149,164,187,.18)' },
                },
            },
        },
    })
}

async function loadComparison() {
    if (!selectedPlayer1 || !selectedPlayer2) {
        setEmpty(comparisonCard, 'Выбери двух игроков для сравнения')
        return
    }

    updateURLParam('player1', selectedPlayer1)
    updateURLParam('player2', selectedPlayer2)
    updateURLParam('season', getSelectedSeason())
    setLoading(comparisonCard, 'Сравниваю игроков...')

    try {
        const data = await apiFetch(`/analytics/compare?player1=${encodeURIComponent(selectedPlayer1)}&player2=${encodeURIComponent(selectedPlayer2)}&season=${encodeURIComponent(getSelectedSeason())}`)
        renderComparison(data)
        renderRadar(data)
    } catch (error) {
        setError(comparisonCard, error.message)
    }
}

function bindSearch(inputId, resultsId, targetId, setter) {
    const input = qs(inputId)
    const results = qs(resultsId)
    input.addEventListener('input', () => renderSearchResults(targetId, input, results, setter))
}

async function init() {
    selectedPlayer1 = params.get('player1')
    selectedPlayer2 = params.get('player2')

    try {
        await loadSeasonOptions(seasonSelect, params.get('season') || '')
        await loadPlayersForSeason()

        if (selectedPlayer1) {
            const player = seasonPlayers.find(item => String(item.player_id) === String(selectedPlayer1))
            if (player) qs('#player1Input').value = player.player
        }

        if (selectedPlayer2) {
            const player = seasonPlayers.find(item => String(item.player_id) === String(selectedPlayer2))
            if (player) qs('#player2Input').value = player.player
        }

        bindSearch('#player1Input', '#player1Results', 'player1Results', value => selectedPlayer1 = value)
        bindSearch('#player2Input', '#player2Results', 'player2Results', value => selectedPlayer2 = value)

        seasonSelect.addEventListener('change', async () => {
            await loadPlayersForSeason()
            await loadComparison()
        })

        qs('#compareButton').addEventListener('click', loadComparison)
        await loadComparison()
    } catch (error) {
        setError(comparisonCard, error.message)
    }
}

init()
