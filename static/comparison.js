let radarChart = null

const API =
    'http://localhost:8080'

const params =
    new URLSearchParams(
        window.location.search
    )

const player1 =
    params.get('player1')

const player2 =
    params.get('player2')

const seasonSelect =
    document.getElementById(
        'seasonSelect'
    )

function getSelectedSeason() {

    return seasonSelect.value
}

async function loadSeasons() {

    const response =
        await fetch(
            `${API}/seasons`
        )

    const seasons =
        await response.json()

    seasons.forEach(season => {

        const option =
            document.createElement(
                'option'
            )

        option.value =
            season

        option.textContent =
            season

        seasonSelect.appendChild(
            option
        )
    })
}

seasonSelect.addEventListener(
    'change',
    loadComparison
)

async function loadComparison() {

    let url =
        `${API}/analytics/compare?player1=${player1}&player2=${player2}`

    const season =
        getSelectedSeason()

    if (season) {

        url += `&season=${season}`
    }

    const response =
        await fetch(url)

    const data =
        await response.json()

    renderComparison(data)

    renderRadar(data)
}

function winnerClass(a, b) {

    return a > b
        ? 'winner'
        : ''
}

function renderComparison(data) {

    document.getElementById(
        'comparisonCard'
    ).innerHTML = `

        <div class="compare-grid">

            <div class="compare-player">

                <h2>
                    ${data.player1.player}
                </h2>

                <p>
                    ${data.player1.team}
                </p>

                <p>
                    ${data.player1.position}
                </p>

            </div>

            <div class="compare-player">

                <h2>
                    ${data.player2.player}
                </h2>

                <p>
                    ${data.player2.team}
                </p>

                <p>
                    ${data.player2.position}
                </p>

            </div>

        </div>

        <table class="compare-table">

            <tr>
                <th>Stat</th>
                <th>${data.player1.player}</th>
                <th>${data.player2.player}</th>
            </tr>

            <tr>
                <td>Goals</td>

                <td class="${winnerClass(
                    data.player1.goals,
                    data.player2.goals
                )}">
                    ${data.player1.goals}
                </td>

                <td class="${winnerClass(
                    data.player2.goals,
                    data.player1.goals
                )}">
                    ${data.player2.goals}
                </td>
            </tr>

            <tr>
                <td>Assists</td>

                <td class="${winnerClass(
                    data.player1.assists,
                    data.player2.assists
                )}">
                    ${data.player1.assists}
                </td>

                <td class="${winnerClass(
                    data.player2.assists,
                    data.player1.assists
                )}">
                    ${data.player2.assists}
                </td>
            </tr>

            <tr>
                <td>Points</td>

                <td class="${winnerClass(
                    data.player1.points,
                    data.player2.points
                )}">
                    ${data.player1.points}
                </td>

                <td class="${winnerClass(
                    data.player2.points,
                    data.player1.points
                )}">
                    ${data.player2.points}
                </td>
            </tr>

            <tr>
                <td>Overall</td>

                <td class="${winnerClass(
                    data.player1.overall_score,
                    data.player2.overall_score
                )}">
                    ${data.player1.overall_score.toFixed(1)}
                </td>

                <td class="${winnerClass(
                    data.player2.overall_score,
                    data.player1.overall_score
                )}">
                    ${data.player2.overall_score.toFixed(1)}
                </td>
            </tr>

        </table>
    `
}

function renderRadar(data) {

    const ctx =
        document.getElementById(
            'comparisonChart'
        )

    const maxGoals =
        Math.max(
            data.player1.goals,
            data.player2.goals
        )

    const maxAssists =
        Math.max(
            data.player1.assists,
            data.player2.assists
        )

    const maxHits =
        Math.max(
            data.player1.hits,
            data.player2.hits
        )

    const maxBlocks =
        Math.max(
            data.player1.blocks,
            data.player2.blocks
        )

    const maxOverall =
        Math.max(
            data.player1.overall_score,
            data.player2.overall_score
        )

    function normalize(
        value,
        max
    ) {

        if (max === 0) {
            return 0
        }

        return (
            value / max
        ) * 100
    }

    if (radarChart) {
        radarChart.destroy()
    }

    radarChart =
        new Chart(ctx, {

            type: 'radar',

            data: {

                labels: [
                    'Goals',
                    'Assists',
                    'Hits',
                    'Blocks',
                    'Overall'
                ],

                datasets: [

                    {
                        label:
                            data.player1.player,

                        data: [

                            normalize(
                                data.player1.goals,
                                maxGoals
                            ),

                            normalize(
                                data.player1.assists,
                                maxAssists
                            ),

                            normalize(
                                data.player1.hits,
                                maxHits
                            ),

                            normalize(
                                data.player1.blocks,
                                maxBlocks
                            ),

                            normalize(
                                data.player1.overall_score,
                                maxOverall
                            )
                        ]
                    },

                    {
                        label:
                            data.player2.player,

                        data: [

                            normalize(
                                data.player2.goals,
                                maxGoals
                            ),

                            normalize(
                                data.player2.assists,
                                maxAssists
                            ),

                            normalize(
                                data.player2.hits,
                                maxHits
                            ),

                            normalize(
                                data.player2.blocks,
                                maxBlocks
                            ),

                            normalize(
                                data.player2.overall_score,
                                maxOverall
                            )
                        ]
                    }
                ]
            },

            options: {

                responsive: true,

                maintainAspectRatio: false,

                scales: {

                    r: {

                        min: 0,

                        max: 100,

                        ticks: {
                            display: false
                        }
                    }
                }
            }
        })
    }

async function init() {

    await loadSeasons()

    await loadComparison()
}

init()