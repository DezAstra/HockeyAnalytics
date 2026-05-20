const API =
    'http://localhost:8080'

const params =
    new URLSearchParams(
        window.location.search
    )

const playerId =
    params.get('id')

console.log("LOAD START")
async function loadPlayer() {

    const response =
        await fetch(
            `${API}/players/${playerId}/career`
        )

    const data =
        await response.json()

    console.log("RENDER")
    console.log(data)
    console.log("CHART")
    console.log(data.career)

    console.log(
        data.career[
            data.career.length - 1
        ]
    )
    
    renderPlayer(data)
    console.log(
        data.career[
            data.career.length - 1
        ]
    )
    renderChart(data.career)
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

loadPlayer()