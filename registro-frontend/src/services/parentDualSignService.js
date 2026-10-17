import api from './api'

export const parentDualSignService = {
    listAuthorizations() {
        return api.get('/parents/dual-authorizations')
    },
    signAuthorization(id, pin) {
        return api.post(`/parents/dual-authorizations/${id}/sign`, { pin })
    },
    rejectAuthorization(id, reason) {
        return api.post(`/parents/dual-authorizations/${id}/reject`, { reason })
    },
    getChildCustodyInfo(studentId) {
        return api.get(`/parents/child/${studentId}/custody`)
    }
}

export default parentDualSignService
