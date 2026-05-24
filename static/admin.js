const csvForm = qs('#csvForm')
const syncForm = qs('#syncForm')
const resultLog = qs('#resultLog')
const warningBox = qs('#warningBox')
const importLogs = qs('#importLogs')

function renderResult(result) {
    resultLog.textContent = JSON.stringify(result, null, 2)

    const errors = result?.result?.errors || result?.errors || []
    const skipped = result?.result?.skipped || result?.skipped || 0

    if (errors.length || skipped > 0) {
        warningBox.innerHTML = `
            <div class="state state-warning">
                Операция завершилась с предупреждениями: skipped=${skipped}, errors=${errors.length}. Проверь JSON ниже.
            </div>
        `
    } else {
        warningBox.innerHTML = '<div class="state state-success">Операция завершилась без предупреждений.</div>'
    }
}

async function loadLogs() {
    setLoading(importLogs, 'Загружаю историю...')
    try {
        const logs = await apiFetch('/import/logs?limit=50')
        if (!logs.length) {
            setEmpty(importLogs, 'История импортов пока пустая')
            return
        }

        importLogs.innerHTML = logs.map(log => `
            <tr>
                <td>${escapeHTML(log.created_at)}</td>
                <td>${escapeHTML(log.source)}</td>
                <td>${escapeHTML(log.season || '—')}</td>
                <td><span class="status-pill status-${escapeHTML(log.status)}">${escapeHTML(log.status)}</span></td>
                <td>${log.imported}</td>
                <td>${log.updated}</td>
                <td>${log.skipped}</td>
                <td>${escapeHTML((log.errors || []).slice(0, 2).join(' | ') || '—')}</td>
            </tr>
        `).join('')
    } catch (error) {
        setError(importLogs, error.message)
    }
}

csvForm?.addEventListener('submit', async event => {
    event.preventDefault()
    const file = qs('#csvFile')?.files?.[0]
    if (!file) return

    const formData = new FormData()
    formData.append('file', file)
    resultLog.textContent = 'Импортирую CSV...'

    try {
        const result = await apiFetch('/import/csv', { method: 'POST', body: formData })
        renderResult(result)
        await loadLogs()
    } catch (error) {
        setError(warningBox, error.message)
        resultLog.textContent = error.message
    }
})

syncForm?.addEventListener('submit', async event => {
    event.preventDefault()
    const season = qs('#syncSeason')?.value?.trim()
    if (!season) return
    resultLog.textContent = 'Синхронизирую сезон через NHL API...'

    try {
        const result = await apiFetch(`/nhl/sync?season=${encodeURIComponent(season)}`, { method: 'POST' })
        renderResult(result)
        await loadLogs()
    } catch (error) {
        setError(warningBox, error.message)
        resultLog.textContent = error.message
    }
})

qs('#refreshLogs')?.addEventListener('click', loadLogs)
loadLogs()
