import { ref } from 'vue';
import { useSchoolStore } from '../stores/schools';
import { useQuasar } from 'quasar';

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
        try {
            if (editingSchool.value) {
                await store.updateSchool(editingSchool.value.id, data);
                $q.notify({ type: 'positive', message: 'School updated' });
            } else {
                await store.createSchool(data);
                $q.notify({ type: 'positive', message: 'School created' });
            }
            showDialog.value = false;
        } catch (e) {
            $q.notify({ type: 'negative', message: 'Operation failed' });
        }
    };

    const confirmDelete = (id) => {
        $q.dialog({
            title: 'Confirm',
            message: 'Are you sure you want to delete this school?',
            cancel: true,
            persistent: true
        }).onOk(async () => {
            await store.deleteSchool(id);
            $q.notify({ type: 'positive', message: 'School deleted' });
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
