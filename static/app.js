const API = 'http://localhost:8080'

const table =
    document.getElementById(
        'playersTable'
    )

const seasonSelect =
    document.getElementById(
        'seasonSelect'
    )

const searchInput =
    document.getElementById(
        'searchInput'
    )

const sortSelect =
    document.getElementById(
        'sortSelect'
    )

let allPlayers = []
function getScoreClass(score) {

    if (score >= 85) {
        return 'score-elite'
    }

    if (score >= 70) {
        return 'score-good'
    }

    if (score >= 55) {
        return 'score-average'
    }

    if (score >= 40) {
        return 'score-bad'
    }

    return 'score-awful'
}

async function loadAnalytics() {

    const season =
        seasonSelect.value

    const response =
        await fetch(
            `${API}/analytics?season=${season}`
        )

    const data =
        await response.json()

    allPlayers = data

    applyFilters()
}

function applyFilters() {

    let players =
        [...allPlayers]

    const search =
        searchInput.value
            .toLowerCase()

    if (search) {

        players =
            players.filter(
                player =>
                    player.player
                        .toLowerCase()
                        .includes(search)
            )
    }

    const sort =
        sortSelect.value

    if (sort === 'overall') {

        players.sort(
            (a, b) =>
                b.overall_score -
                a.overall_score
        )
    }

    if (sort === 'normalized') {

        players.sort(
            (a, b) =>
                b.normalized_score -
                a.normalized_score
        )
    }

    if (sort === 'context') {

        players.sort(
            (a, b) =>
                b.context_score -
                a.context_score
        )
    }

    if (sort === 'percentile') {

        players.sort(
            (a, b) =>
                b.overall_percentile -
                a.overall_percentile
        )
    }

    renderPlayers(players)
}

function renderPlayers(players) {

    table.innerHTML = ''

    players.forEach(player => {

        const row =
            document.createElement('tr')

        row.innerHTML = `
            <td>
                <a href="/player?id=${player.player_id}">
                    ${player.player}
                </a>
            </td>

            <td>
                ${player.team}
            </td>

            <td>
                ${player.position}
            </td>

            <td class="${getScoreClass(
                player.normalized_score
            )}">
                ${player.normalized_score.toFixed(1)}
            </td>

            <td class="${getScoreClass(
                player.context_score
            )}">
                ${player.context_score.toFixed(1)}
            </td>

            <td class="${getScoreClass(
                player.overall_score
            )}">
                ${player.overall_score.toFixed(1)}
            </td>

            <td>
                ${player.overall_percentile.toFixed(1)}
            </td>

            <td>
                <span class="badge">
                    ${player.archetype}
                </span>
            </td>
        `

        table.appendChild(row)
    })
}

seasonSelect.addEventListener(
    'change',
    loadAnalytics,
)

searchInput.addEventListener(
    'input',
    applyFilters,
)

sortSelect.addEventListener(
    'change',
    applyFilters,
)

loadAnalytics()