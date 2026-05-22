const API =
    'http://localhost:8080'

const params =
    new URLSearchParams(
        window.location.search
    )

const playerId =
    params.get('id')

async function loadPlayer() {

    const response =
        await fetch(
            `${API}/players/${playerId}/career`
        )

    const data =
        await response.json()

    renderPlayer(data)

    await loadSimilarPlayers()
    console.log(
        data.career[
            data.career.length - 1
        ]
    )
    renderChart(data.career)
}

async function loadSimilarPlayers() {

    const response =
        await fetch(
            `${API}/players/${playerId}/similar`
        )

    const data =
        await response.json()

    renderSimilarPlayers(data)
}

function renderSimilarPlayers(players) {

    document.getElementById(
        'similarPlayers'
    ).innerHTML = `

        <h2>
            Similar Players
        </h2>

        <div class="similar-grid">

            ${players.map(player => `

                <a
                    href="/player?id=${player.player_id}"
                    class="similar-card"
                >

                    <h3>
                        ${player.player}
                    </h3>

                    <p>
                        ${player.team}
                    </p>

                    <p>
                        ${player.archetype}
                    </p>

                    <div class="similarity-score">

                        ${player.similarity.toFixed(1)}%
                        similar

                    </div>

                </a>

            `).join('')}

        </div>
    `
}

function renderPlayer(data) {

    const latest =
        data.career[
            data.career.length - 1
        ]

    document.getElementById(
        'playerCard'
    ).innerHTML = `

        <h1 style="
            font-size:42px;
            margin-bottom:10px;
        ">
            ${data.player}
        </h1>

        <div style="
            color:#94a3b8;
            margin-bottom:20px;
        ">
            ${data.position}
        </div>

        <span class="badge">
            ${latest.archetype}
        </span>

        <div class="stat-grid">

            <div class="stat-box">

                <div class="stat-title">
                    Goals
                </div>

                <div class="stat-value">
                    ${latest.goals}
                </div>

            </div>

            <div class="stat-box">

                <div class="stat-title">
                    Assists
                </div>

                <div class="stat-value">
                    ${latest.assists}
                </div>

            </div>

            <div class="stat-box">

                <div class="stat-title">
                    Overall
                </div>

                <div class="stat-value">
                    ${latest.overall_score.toFixed(1)}
                </div>

            </div>

            <div class="stat-box">

                <div class="stat-title">
                    Games
                </div>

                <div class="stat-value">
                    ${latest.games_played}
                </div>

            </div>

        </div>
    `
}

function renderChart(career) {

    const ctx =
        document
            .getElementById(
                'careerChart'
            )

    new Chart(ctx, {
        type: 'line',

        data: {
            labels:
                career.map(
                    s => s.season
                ),

            datasets: [{
                label:
                    'Overall Score',

                data:
                    career.map(
                        s =>
                            s.overall_score
                    ),

                borderColor:
                    '#67e8f9',

                tension: 0.3
            }]
        }
    })
}

async function searchPlayers(query) {

    const response =
        await fetch(
            `${API}/analytics?season=25/26`
        )

    const players =
        await response.json()

    return players.filter(player => {

        if (
            player.player_id ==
            playerId
        ) {
            return false
        }

        return player.player
            .toLowerCase()
            .includes(
                query.toLowerCase()
            )
    })
}

document
    .getElementById(
        'compareInput'
    )
    .addEventListener(
        'input',
        async e => {

            const query =
                e.target.value

            const results =
                document.getElementById(
                    'compareResults'
                )

            if (
                query.length < 2
            ) {

                results.innerHTML = ''

                return
            }

            const players =
                await searchPlayers(
                    query
                )

            results.innerHTML =
                players
                    .slice(0, 8)
                    .map(player => `

                        <div
                            class="compare-result"
                            data-id="${player.player_id}"
                        >

                            ${player.player}
                            (${player.team})

                        </div>

                    `)
                    .join('')

            document
                .querySelectorAll(
                    '.compare-result'
                )
                .forEach(card => {

                    card.addEventListener(
                        'click',
                        () => {

                            const target =
                                card.dataset.id

                            window.location.href =
                                `/comparison?player1=${playerId}&player2=${target}`
                        }
                    )
                })
        }
    )

async function init() {

    await loadPlayer()
}

init()