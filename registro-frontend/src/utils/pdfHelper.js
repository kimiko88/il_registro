import api from '../services/api'

/**
 * Handles PDF download for both synchronous (direct PDF blob) and asynchronous (202 Accepted job queue) responses.
 * @param {Function} fetchFn Function returning an Axios promise with responseType: 'blob'
 * @param {String} defaultFilename Filename for the downloaded PDF
 * @param {Object} options Optional options like pollIntervalMs (default 1500), maxPollAttempts (default 40)
 */
export async function handleAsyncPdfDownload(fetchFn, defaultFilename = 'documento.pdf', options = {}) {
    const pollIntervalMs = options.pollIntervalMs || 1500
    const maxPollAttempts = options.maxPollAttempts || 40

    const response = await fetchFn()
    const contentType = response.headers?.['content-type'] || ''

    let blobData = response.data
    let isAsyncJob = false
    let jobInfo = null

    // If HTTP status is 202 or contentType includes application/json
    if (response.status === 202 || contentType.includes('application/json')) {
        isAsyncJob = true
    } else if (blobData instanceof Blob && blobData.type && blobData.type.includes('application/json')) {
        isAsyncJob = true
    }

    if (isAsyncJob) {
        try {
            const text = blobData instanceof Blob ? await blobData.text() : (typeof blobData === 'string' ? blobData : JSON.stringify(blobData))
            jobInfo = typeof text === 'string' ? JSON.parse(text) : text
        } catch (e) {
            console.warn('Failed to parse async job JSON response from blob:', e)
        }
    }

    if (jobInfo && (jobInfo.job_id || jobInfo.status_url)) {
        const statusUrl = jobInfo.status_url || `/scrutiny/pdf-jobs/${jobInfo.job_id}`
        let attempts = 0
        let completedJob = null

        while (attempts < maxPollAttempts) {
            await new Promise(resolve => setTimeout(resolve, pollIntervalMs))
            attempts++

            const statusRes = await api.get(statusUrl)
            const statusData = statusRes.data

            if (statusData.status === 'completed') {
                completedJob = statusData
                break
            } else if (statusData.status === 'failed') {
                throw new Error(statusData.error || 'Generazione PDF fallita nel worker di background')
            }
        }

        if (!completedJob) {
            throw new Error('Timeout durante la generazione asincrona del PDF')
        }

        // Fetch completed PDF blob (pass sync=true to get actual bytes)
        const rawUrl = completedJob.download_url || statusUrl
        const downloadUrl = rawUrl.includes('?') ? `${rawUrl}&sync=true` : `${rawUrl}?sync=true`
        const finalPdfRes = await api.get(downloadUrl, { responseType: 'blob', timeout: 60000 })
        blobData = finalPdfRes.data
    }

    // Trigger browser download
    const blob = blobData instanceof Blob ? blobData : new Blob([blobData], { type: 'application/pdf' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', defaultFilename)
    document.body.appendChild(link)
    link.click()

    setTimeout(() => {
        if (link.parentNode) link.parentNode.removeChild(link)
        window.URL.revokeObjectURL(url)
    }, 1000)

    return true
}

export default handleAsyncPdfDownload
