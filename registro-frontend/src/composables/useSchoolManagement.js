import { ref } from 'vue';
import { useSchoolStore } from '../stores/schools';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useSchoolManagement() {
    const store = useSchoolStore();
    const $q = useQuasar();

    const showDialog = ref(false);
    const editingSchool = ref(null);

    const onRequest = (props) => {
        store.fetchSchools(props.pagination);
    };

    const openCreate = () => {
        editingSchool.value = null;
        showDialog.value = true;
    };

    const openEdit = (school) => {
        editingSchool.value = { ...school };
        showDialog.value = true;
    };

    const submitSchool = async (data) => {
        const t = i18n?.global?.t;
        try {
            if (editingSchool.value) {
                await store.updateSchool(editingSchool.value.id, data);
                $q.notify({ type: 'positive', message: t ? t('composables.schools.updated') : 'Scuola aggiornata con successo' });
            } else {
                await store.createSchool(data);
                $q.notify({ type: 'positive', message: t ? t('composables.schools.created') : 'Scuola creata con successo' });
            }
            showDialog.value = false;
        } catch (e) {
            $q.notify({ type: 'negative', message: t ? t('composables.schools.operationFailed') : 'Operazione fallita' });
        }
    };

    const confirmDelete = (id) => {
        const t = i18n?.global?.t;
        $q.dialog({
            title: t ? t('composables.schools.deleteConfirmTitle') : 'Conferma eliminazione',
            message: t ? t('composables.schools.deleteConfirmMsg') : 'Sei sicuro di voler eliminare questa scuola?',
            cancel: true,
            persistent: true
        }).onOk(async () => {
            await store.deleteSchool(id);
            $q.notify({ type: 'positive', message: t ? t('composables.schools.deleted') : 'Scuola eliminata con successo' });
        });
    };

    return {
        store,
        showDialog,
        editingSchool,
        onRequest,
        openCreate,
        openEdit,
        submitSchool,
        confirmDelete
    };
}
