import { computed } from 'vue';
import { useChildrenStore } from 'src/stores/children';
import { useQuasar } from 'quasar';

export function useChildrenManagement() {
    const store = useChildrenStore();
    const $q = useQuasar();

    const addChild = async (_childData) => {
        // Mock Add
        $q.notify({ type: 'positive', message: 'Child added successfully' });
        // In real app, would call store action
    };

    const removeChild = async (_childId) => {
        $q.dialog({
            title: 'Confirm',
            message: 'Are you sure you want to remove this child from your profile?',
            cancel: true,
            persistent: true
        }).onOk(() => {
            $q.notify({ type: 'positive', message: 'Child removed' });
        });
    };

    return {
        children: computed(() => store.children),
        addChild,
        removeChild
    };
}
