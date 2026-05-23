const API = 'http://localhost:8080'
const params = new URLSearchParams(window.location.search)
const playerId = params.get('id')
const seasonSelect = document.getElementById('seasonSelect')

// Кэш для данных игрока, чтобы не перезапрашивать `/career` при каждом изменении сезона
let playerDataCache = null

async function loadSeasons(playerCareer) {
    const response = await fetch(`${API}/seasons`)
    const allSeasons = await response.json()

    // Сортируем сезоны по убыванию (хронологически новые сверху)
    allSeasons.sort((a, b) => {
        const aStart = parseInt(a.split('/')[0])
        const bStart = parseInt(b.split('/')[0])
        return bStart - aStart
    })

    // Создаем Set из сезонов, в которых игрок реально играл
    const playedSeasons = new Set(playerCareer.map(s => s.season))

    seasonSelect.innerHTML = ''

    allSeasons.forEach(season => {
        // ДОБАВЛЕНА ПРОВЕРКА: Проверяем, был ли игрок в этом сезоне
        if (playedSeasons.has(season)) {
            const option = document.createElement('option')
            option.value = season
            option.textContent = `Season ${season}`
            seasonSelect.appendChild(option)
        }
    })

    // Выставляем дефолтное значение: самый свежий сезон, в котором играл этот конкретный игрок
    if (seasonSelect.options.length > 0) {
        seasonSelect.value = seasonSelect.options[0].value
    }
}

function getSelectedSeason() {
    return seasonSelect.value
}

async function loadPlayer() {
    // Если данные игрока ещё не загружены — загружаем один раз
    if (!playerDataCache) {
        const response = await fetch(`${API}/players/${playerId}/career`)
        playerDataCache = await response.json()
    }

    // Отрисовываем профиль на основе выбранного сезона
    renderPlayer(playerDataCache)
    
    // Загружаем похожих игроков именно под этот сезон
    await loadSimilarPlayers()
    
    // Строим график по всей карьере
    renderChart(playerDataCache.career)
}

async function loadSimilarPlayers() {
    const currentSeason = getSelectedSeason()
    if (!currentSeason) return

    const response = await fetch(`${API}/players/${playerId}/similar?season=${currentSeason}`)
    const data = await response.json()
    renderSimilarPlayers(data)
}

function renderSimilarPlayers(players) {
    document.getElementById('similarPlayers').innerHTML = `
        <h2>Similar Players</h2>
        <div class="similar-grid">
            ${players.map(player => `
                <a href="/player?id=${player.player_id}" class="similar-card">
                    <h3>${player.player}</h3>
                    <p>${player.team}</p>
                    <p>${player.archetype}</p>
                    <div class="similarity-score">
                        ${player.similarity.toFixed(1)}% similar
                    </div>
                </a>
            `).join('')}
        </div>
    `
}

function renderPlayer(data) {
    const selectedSeason = getSelectedSeason()

    const latest = data.career.find(season => season.season === selectedSeason) || data.career[data.career.length - 1]

    // Определяем правильный CSS класс для цвета
    const archetypeClass = getArchetypeClass(latest.archetype)

    document.getElementById('playerCard').innerHTML = `
        <h1 style="font-size:42px; margin-bottom:10px;">
            ${data.player}
        </h1>
        <div style="color:#94a3b8; margin-bottom:20px;">
            ${data.position}
        </div>
        
        <span class="badge ${archetypeClass}">
            ${latest.archetype || 'Баланс'}
        </span>

        <div class="stat-grid" style="margin-top: 15px;"> <div class="stat-box">
                <div class="stat-title">Goals</div>
                <div class="stat-value">${latest.goals}</div>
            </div>
            <div class="stat-box">
                <div class="stat-title">Assists</div>
                <div class="stat-value">${latest.assists}</div>
            </div>
            <div class="stat-box">
                <div class="stat-title">Overall</div>
                <div class="stat-value">${latest.overall_score.toFixed(1)}</div>
            </div>
            <div class="stat-box">
                <div class="stat-title">Games</div>
                <div class="stat-value">${latest.games_played}</div>
            </div>
        </div>
    `
}

let careerChart = null

function renderChart(career) {
    const ctx = document.getElementById('careerChart')
    if (careerChart) {
        careerChart.destroy()
    }

    // Хронологический порядок для графика (от старых к новым)
    const sortedCareer = [...career].sort((a, b) => {
        return parseInt(a.season.split('/')[0]) - parseInt(b.season.split('/')[0])
    })

    careerChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: sortedCareer.map(s => s.season),
            datasets: [{
                label: 'Overall Score',
                data: sortedCareer.map(s => s.overall_score),
                borderColor: '#67e8f9',
                tension: 0.3
            }]
        }
    })
}

async function searchPlayers(query) {
    const response = await fetch(`${API}/analytics?season=${getSelectedSeason()}`)
    const players = await response.json()

    return players.filter(player => {
        if (player.player_id == playerId) {
            return false
        }
        return player.player.toLowerCase().includes(query.toLowerCase())
    })
}

document.getElementById('compareInput').addEventListener('input', async e => {
    const query = e.target.value
    const results = document.getElementById('compareResults')

    if (query.length < 2) {
        results.innerHTML = ''
        return
    }

    const players = await searchPlayers(query)

    results.innerHTML = players
        .slice(0, 8)
        .map(player => `
            <div class="compare-result" data-id="${player.player_id}">
                ${player.player} (${player.team})
            </div>
        `)
        .join('')

    document.querySelectorAll('.compare-result').forEach(card => {
        card.addEventListener('click', () => {
            const target = card.dataset.id
            window.location.href = `/comparison?player1=${playerId}&player2=${target}`
        })
    })
})

// При переключении сезона не качаем заново всю карьеру, а берем из кэша
seasonSelect.addEventListener('change', async () => {
    if (playerDataCache) {
        renderPlayer(playerDataCache)
        await loadSimilarPlayers()
    }
})

// Изменен порядок инициализации
async function init() {
    if (!playerId) {
        console.error('No player ID provided in URL')
        return
    }

    // 1. Сначала загружаем профиль и карьеру игрока
    const response = await fetch(`${API}/players/${playerId}/career`)
    playerDataCache = await response.json()

    // 2. На основе карьеры строим проверенный список сезонов
    await loadSeasons(playerDataCache.career)

    // 3. Рендерим интерфейс
    renderPlayer(playerDataCache)
    await loadSimilarPlayers()
    renderChart(playerDataCache.career)
}

function getArchetypeClass(archetype) {
    if (!archetype) return 'badge-default';
    
    // Приводим к одному регистру и убираем пробелы на всякий случай
    const arch = archetype.trim();

    switch (arch) {
        case 'Снайпер': return 'badge-sniper';
        case 'Бомбардир': return 'badge-pointer';
        case 'Ассистент': return 'badge-assistman';
        case 'Атакующий защитник': return 'badge-offensive-defenseman';
        case 'Защитник-стена': return 'badge-iron-defenseman';
        case 'Нарушитель': return 'badge-offender';
        case 'Силовик': return 'badge-grinder';
        case 'Специалист по вбрасываниям': return 'badge-faceoff-specialist';
        default: return 'badge-default';
    }
}

init()