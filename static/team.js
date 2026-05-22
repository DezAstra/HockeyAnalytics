const API =
    'http://localhost:8080'

let currentSort = null

let sortDirection = 'desc'

const params =
    window.location.pathname.split('/')

const team =
    params[
        params.length - 1
    ]

const table =
    document.getElementById(
        'playersTable'
    )

const seasonSelect =
    document.getElementById(
        'seasonSelect'
    )

async function loadTeam() {

    const season =
        seasonSelect.value

    const response =
        await fetch(
            `${API}/api/team/${team}?season=${season}`
        )

    const data =
        await response.json()

    renderTeam(data)
}

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

function renderTeam(data) {

        let players =
        [...data.players]

        if (currentSort) {

            players.sort((a, b) => {

                const aVal =
                    a[currentSort]

                const bVal =
                    b[currentSort]

                if (
                    sortDirection === 'asc'
                ) {

                    return aVal - bVal
                }

                return bVal - aVal
            })
        }

    document.getElementById(
        'teamHeader'
    ).textContent =
        `${data.team} — ${data.season}`

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

                <span class="badge ${getArchetypeClass(player.archetype)}">

                    ${player.archetype}

                </span>

            </td>
        `

        table.appendChild(row)
    })
}

seasonSelect.addEventListener(
    'change',
    loadTeam
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
                        sortDirection === 'asc'
                            ? 'desc'
                            : 'asc'

                } else {

                    currentSort = field

                    sortDirection = 'desc'
                }

                loadTeam()
            }
        )
    })

loadTeam()