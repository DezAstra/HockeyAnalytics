const API = 'http://localhost:8080'

let allPlayers = []

let currentSort = null

let sortDirection = 'desc'

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

    populateTeamFilter(data)

    renderPlayers(data)
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
                    <a href="/team/${player.team}">
                        ${player.team}
                    </a>
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
                <span class="badge ${getArchetypeClass(player.archetype)}">
                    ${player.archetype}
                </span>
            </td>
        `

        table.appendChild(row)
    })
}

function applyFilters() {

    let filtered =
        [...allPlayers]

    const search =
        document
            .getElementById(
                'searchInput'
            )
            .value
            .toLowerCase()

    const position =
        document
            .getElementById(
                'positionFilter'
            )
            .value

    const team =
        document
            .getElementById(
                'teamFilter'
            )
            .value

    const archetype =
        document
            .getElementById(
                'archetypeFilter'
            )
            .value

    if (search) {

        filtered =
            filtered.filter(
                p =>
                    p.player
                        .toLowerCase()
                        .includes(search)
            )
    }

    if (position) {

        filtered =
            filtered.filter(
                p =>
                    p.position ===
                    position
            )
    }

    if (team) {

        filtered =
            filtered.filter(
                p =>
                    p.team ===
                    team
            )
    }

    if (archetype) {

        filtered =
            filtered.filter(
                p =>
                    p.archetype ===
                    archetype
            )
    }

    if (currentSort) {

        filtered.sort(
            (a, b) => {

                const aVal =
                    a[currentSort]

                const bVal =
                    b[currentSort]

                if (
                    sortDirection ===
                    'asc'
                ) {

                    return aVal - bVal
                }

                return bVal - aVal
            }
        )
    }

    renderPlayers(filtered)
}

function getArchetypeClass(archetype) {

    switch (archetype) {

        case 'Ассистент':
            return 'badge-assistman'

        case 'Снайпер':
            return 'badge-sniper'

        case 'Бомбардир':
            return 'badge-pointer'

        case 'Защитник-стена':
            return 'badge-iron-defenseman'

        case 'Атакующий защитник':
            return 'badge-offensive-defenseman'

        case 'Нарушитель':
            return 'badge-offender'

        case 'Силовик':
            return 'badge-grinder'

        case 'Специалист по вбрасываниям':
            return 'badge-faceoff-specialist'

        default:
            return ''
    }
}

function populateTeamFilter(data) {

    const select =
        document.getElementById(
            'teamFilter'
        )

    select.innerHTML =
        `
        <option value="">
            Все команды
        </option>
        `

    const teams =
        [...new Set(
            data.map(
                p => p.team
            )
        )]

    teams.sort()

    teams.forEach(team => {

        const option =
            document.createElement(
                'option'
            )

        option.value =
            team

        option.textContent =
            team

        select.appendChild(option)
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

document
    .getElementById(
        'searchInput'
    )
    .addEventListener(
        'input',
        applyFilters
    )

document
    .getElementById(
        'positionFilter'
    )
    .addEventListener(
        'change',
        applyFilters
    )

document
    .getElementById(
        'teamFilter'
    )
    .addEventListener(
        'change',
        applyFilters
    )

document
    .getElementById(
        'archetypeFilter'
    )
    .addEventListener(
        'change',
        applyFilters
    )

document
    .querySelectorAll(
        'th[data-sort]'
    )
    .forEach(th => {

        th.addEventListener(
            'click',
            () => {

                const field =
                    th.dataset.sort

                if (
                    currentSort === field
                ) {

                    sortDirection =
                        sortDirection ===
                        'asc'
                            ? 'desc'
                            : 'asc'

                } else {

                    currentSort =
                        field

                    sortDirection =
                        'desc'
                }

                applyFilters()
            }
        )
    })

loadAnalytics()