const params = new URLSearchParams(window.location.search)
const playerId = params.get('id')
const seasonSelect = qs('#seasonSelect')
let playerDataCache = null
let careerChart = null
let analyticsCache = new Map()

function getSelectedSeason() {
    return seasonSelect.value
}


function getSortedCareer(career) {
    return sortSeasonsAsc((career || []).map(item => item.season))
        .map(season => career.find(item => item.season === season))
        .filter(Boolean)
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

function getTrendInfo(delta) {
    const value = Number(delta)

    if (!Number.isFinite(value)) {
        return {
            label: 'Недостаточно данных',
            className: 'trend-flat',
            icon: '—',
        }
    }

    if (value > 1) {
        return {
            label: 'Прогрессирует',
            className: 'trend-up',
            icon: '↗',
        }
    }

    if (value < -1) {
        return {
            label: 'Проседает',
            className: 'trend-down',
            icon: '↘',
        }
    }

    return {
        label: 'Стабильно',
        className: 'trend-flat',
        icon: '→',
    }
}

function renderProgressionBlock(career, selectedSeason) {
    const sortedCareer = getSortedCareer(career)

    if (sortedCareer.length < 2) {
        return `
            <div class="progression-panel">
                <div class="progression-title">Прогрессия игрока</div>
                <div class="subtitle">Нужно минимум два сезона, чтобы оценить динамику.</div>
            </div>
        `
    }

    const currentIndex = Math.max(
        0,
        sortedCareer.findIndex(item => item.season === selectedSeason),
    )

    const current = sortedCareer[currentIndex] || sortedCareer.at(-1)
    const previous = sortedCareer[currentIndex - 1]
    const first = sortedCareer[0]
    const best = [...sortedCareer].sort(
        (a, b) =>
            Number(b.overall_score || 0) -
            Number(a.overall_score || 0),
    )[0]

    const previousDelta = previous
        ? Number(current.overall_score || 0) - Number(previous.overall_score || 0)
        : null

    const careerDelta = Number(current.overall_score || 0) - Number(first.overall_score || 0)
    const trend = getTrendInfo(previousDelta)

    return `
        <div class="progression-panel">
            <div class="progression-main ${trend.className}">
                <span class="progression-icon">${trend.icon}</span>
                <div>
                    <div class="progression-label">${trend.label}</div>
                    <div class="subtitle">
                        ${previous
                            ? `к прошлому сезону: ${signedNumber(previousDelta)} overall`
                            : 'нет предыдущего сезона для сравнения'}
                    </div>
                </div>
            </div>

            <div class="progression-grid">
                <div class="progression-item">
                    <div class="stat-title">Текущий сезон</div>
                    <div class="stat-value ${getScoreClass(current.overall_score)}">
                        ${formatNumber(current.overall_score)}
                    </div>
                    <div class="subtitle">${escapeHTML(current.season)}</div>
                </div>

                <div class="progression-item">
                    <div class="stat-title">За карьеру</div>
                    <div class="stat-value ${careerDelta >= 0 ? 'trend-up-text' : 'trend-down-text'}">
                        ${signedNumber(careerDelta)}
                    </div>
                    <div class="subtitle">с ${escapeHTML(first.season)}</div>
                </div>

                <div class="progression-item">
                    <div class="stat-title">Лучший сезон</div>
                    <div class="stat-value ${getScoreClass(best.overall_score)}">
                        ${formatNumber(best.overall_score)}
                    </div>
                    <div class="subtitle">${escapeHTML(best.season)}</div>
                </div>
            </div>
        </div>
    `
}

function renderRatingExplanation(season) {
    const strong = []
    const weak = []

    if (Number(season.overall_percentile || 0) >= 85) {
        strong.push('входит в верхнюю группу сезона по «Лучше %»')
    } else if (Number(season.overall_percentile || 0) < 40) {
        weak.push('ниже среднего уровня сезона по «Лучше %»')
    }

    if (Number(season.normalized_score || 0) >= 70) {
        strong.push('хорошая эффективность после нормализации')
    } else if (Number(season.normalized_score || 0) < 40) {
        weak.push('низкая нормализованная эффективность')
    }

    if (Number(season.context_score || 0) >= 70) {
        strong.push('хороший вклад с учётом роли и позиции')
    } else if (Number(season.context_score || 0) < 40) {
        weak.push('контекстная модель оценивает вклад ниже среднего')
    }

    if (Number(season.points || 0) >= 70) strong.push('высокая результативность по очкам')
    if (Number(season.games_played || 0) < 20) weak.push('маленькая выборка матчей может искажать оценку')

    if (!strong.length) strong.push('показатели близки к среднему профилю сезона')
    if (!weak.length) weak.push('явных слабых сигналов по текущим метрикам нет')

    return `
        <div class="explain-panel">
            <h3>Почему такой рейтинг?</h3>
            <div class="grid-2">
                <div>
                    <div class="stat-title">Сильные стороны</div>
                    <ul class="insight-list good">${strong.map(item => `<li>${escapeHTML(item)}</li>`).join('')}</ul>
                </div>
                <div>
                    <div class="stat-title">Риски / слабые стороны</div>
                    <ul class="insight-list warn">${weak.map(item => `<li>${escapeHTML(item)}</li>`).join('')}</ul>
                </div>
            </div>
        </div>
    `
}

function renderPlayer(data) {
    const selectedSeason = getSelectedSeason()
    const latest = data.career.find(item => item.season === selectedSeason) || data.career.at(-1)

    if (!latest) {
        setEmpty(qs('#playerCard'), 'У игрока нет сезонной статистики')
        return
    }

    qs('#playerCard').innerHTML = `
        <div class="player-hero">
            ${playerPhotoHTML(data.nhl_id, data.player, 'player-photo-large')}
            <div class="player-hero-main">
                <div class="title">${escapeHTML(data.player)}</div>
                <div class="subtitle player-meta">
                    <span>${escapeHTML(data.position || 'Позиция не указана')}</span>
                    <span class="dot-separator">·</span>
                    <span class="team-cell inline-team">
                        ${teamLogoHTML(latest.team, 'team-logo-small')}
                        ${escapeHTML(latest.team || '—')}
                    </span>
                    <span class="dot-separator">·</span>
                    <span>сезон ${escapeHTML(latest.season)}</span>
                </div>
            </div>
            <span class="badge ${getArchetypeClass(latest.archetype)}">${escapeHTML(latest.archetype || 'Баланс')}</span>
        </div>
        <div class="stat-grid">
            <div class="stat-box"><div class="stat-title">Матчи</div><div class="stat-value">${latest.games_played}</div></div>
            <div class="stat-box"><div class="stat-title">Очки</div><div class="stat-value">${latest.points}</div></div>
            <div class="stat-box"><div class="stat-title">Голы</div><div class="stat-value">${latest.goals}</div></div>
            <div class="stat-box"><div class="stat-title">Ассисты</div><div class="stat-value">${latest.assists}</div></div>
            <div class="stat-box"><div class="stat-title">Normalized</div><div class="stat-value ${getScoreClass(latest.normalized_score)}">${formatNumber(latest.normalized_score)}</div></div>
            <div class="stat-box"><div class="stat-title">Context</div><div class="stat-value ${getScoreClass(latest.context_score)}">${formatNumber(latest.context_score)}</div></div>
            <div class="stat-box"><div class="stat-title">Overall</div><div class="stat-value ${getScoreClass(latest.overall_score)}">${formatNumber(latest.overall_score)}</div></div>
            <div class="stat-box"><div class="stat-title">Лучше %</div><div class="stat-value">${formatNumber(latest.overall_percentile)}</div></div>
        </div>

        ${renderProgressionBlock(data.career, latest.season)}
        ${renderRatingExplanation(latest)}
    `
}

function renderSimilarPlayers(players) {
    const container = qs('#similarPlayers')

    if (!players.length) {
        setEmpty(container, 'Похожих игроков для этого сезона нет')
        return
    }

    container.innerHTML = `
        <div class="similar-grid">
            ${players.map(player => `
                <a href="/player?id=${encodeURIComponent(player.player_id)}" class="similar-card">
                    ${playerPhotoHTML(player.nhl_id, player.player, 'player-photo-small')}
                    <h3>${escapeHTML(player.player)}</h3>
                    <p class="subtitle">${escapeHTML(player.team || '—')} · ${escapeHTML(player.position || '—')}</p>
                    <span class="badge ${getArchetypeClass(player.archetype)}">${escapeHTML(player.archetype || 'Баланс')}</span>
                    <div class="similarity-score">${formatNumber(player.similarity)}% похожести</div>
                </a>
            `).join('')}
        </div>
    `
}

async function loadSimilarPlayers() {
    const season = getSelectedSeason()
    const container = qs('#similarPlayers')
    if (!season) return

    setLoading(container, 'Ищу похожих игроков...')

    try {
        const data = await apiFetch(`/players/${encodeURIComponent(playerId)}/similar?season=${encodeURIComponent(season)}`)
        renderSimilarPlayers(data)
    } catch (error) {
        setError(container, error.message)
    }
}

function renderChart(career) {
    const ctx = qs('#careerChart')
    const sortedCareer = getSortedCareer(career)

    if (careerChart) {
        careerChart.destroy()
    }

    careerChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: sortedCareer.map(item => item.season),
            datasets: [
                {
                    label: 'Overall',
                    data: sortedCareer.map(item => Number(item.overall_score || 0)),
                    tension: 0.3,
                },
                {
                    label: 'Context',
                    data: sortedCareer.map(item => Number(item.context_score || 0)),
                    tension: 0.3,
                },
                {
                    label: 'Normalized',
                    data: sortedCareer.map(item => Number(item.normalized_score || 0)),
                    tension: 0.3,
                },
                {
                    label: 'Очки',
                    data: sortedCareer.map(item => Number(item.points || 0)),
                    tension: 0.3,
                },
            ],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { labels: { color: '#edf5ff' } },
            },
            scales: {
                x: { ticks: { color: '#95a4bb' }, grid: { color: 'rgba(149,164,187,.12)' } },
                y: { ticks: { color: '#95a4bb' }, grid: { color: 'rgba(149,164,187,.12)' } },
            },
        },
    })
}

async function loadSeasonList(career) {
    const seasons = sortSeasonsDesc(career.map(item => item.season))
    seasonSelect.innerHTML = ''

    seasons.forEach(season => {
        const option = document.createElement('option')
        option.value = season
        option.textContent = season
        seasonSelect.appendChild(option)
    })

    const selected = params.get('season')
    if (selected && seasons.includes(selected)) {
        seasonSelect.value = selected
    }
}

async function searchPlayers(query) {
    const season = getSelectedSeason()

    if (!analyticsCache.has(season)) {
        analyticsCache.set(season, await apiFetch(`/analytics?season=${encodeURIComponent(season)}`))
    }

    return analyticsCache.get(season)
        .filter(player => String(player.player_id) !== String(playerId))
        .filter(player => player.player.toLowerCase().includes(query.toLowerCase()))
        .slice(0, 8)
}

function bindEvents() {
    seasonSelect.addEventListener('change', async () => {
        updateURLParam('season', getSelectedSeason())
        renderPlayer(playerDataCache)
        renderChart(playerDataCache.career)
        await loadSimilarPlayers()
    })

    qs('#compareInput').addEventListener('input', async event => {
        const query = event.target.value.trim()
        const results = qs('#compareResults')

        if (query.length < 2) {
            results.innerHTML = ''
            return
        }

        try {
            const players = await searchPlayers(query)
            results.innerHTML = players.map(player => `
                <div class="compare-result" data-id="${player.player_id}">
                    ${escapeHTML(player.player)} (${escapeHTML(player.team || '—')})
                </div>
            `).join('')

            qsa('.compare-result').forEach(card => {
                card.addEventListener('click', () => {
                    window.location.href = `/comparison?player1=${encodeURIComponent(playerId)}&player2=${encodeURIComponent(card.dataset.id)}&season=${encodeURIComponent(getSelectedSeason())}`
                })
            })
        } catch (error) {
            setError(results, error.message)
        }
    })
}

async function init() {
    if (!playerId) {
        setError(qs('#playerCard'), 'В URL не передан id игрока')
        return
    }

    bindEvents()
    setLoading(qs('#playerCard'), 'Загружаю карточку игрока...')

    try {
        playerDataCache = await apiFetch(`/players/${encodeURIComponent(playerId)}/career`)
        await loadSeasonList(playerDataCache.career)
        renderPlayer(playerDataCache)
        renderChart(playerDataCache.career)
        await loadSimilarPlayers()
    } catch (error) {
        setError(qs('#playerCard'), error.message)
    }
}

init()
