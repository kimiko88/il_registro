<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none">{{ $t('groupsPage.title') }}</h1>
        <div class="text-subtitle2 text-grey-7">{{ $t('groupsPage.subtitle') }}</div>
      </div>
      <q-space />
      <q-btn color="primary" icon="add" :label="$t('groupsPage.newGroup')" @click="openCreateDialog" />
    </div>

    <q-card v-if="loading" class="q-pa-lg text-center shadow-2">
      <q-spinner color="primary" size="3em" />
      <div class="q-mt-sm text-grey-7">{{ $t('groupsPage.loading') }}</div>
    </q-card>

    <div v-else-if="groups.length === 0" class="q-pa-xl text-center">
      <q-icon name="groups" size="4rem" color="grey-5" />
      <div class="text-h6 text-grey-6 q-mt-md">{{ $t('groupsPage.noGroups') }}</div>
      <div class="text-caption text-grey-5 q-mb-md">{{ $t('groupsPage.noGroupsDesc') }}</div>
      <q-btn color="primary" icon="add" :label="$t('groupsPage.createGroup')" @click="openCreateDialog" />
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-for="g in groups" :key="g.id" class="col-12 col-md-6 col-lg-4">
        <q-card class="shadow-2 rounded-borders hover:shadow-4 transition-all">
          <q-card-section class="bg-primary text-white">
            <div class="row items-center no-wrap">
              <div class="col">
                <div class="text-h6 text-weight-bold">{{ g.name }}</div>
                <div class="text-caption opacity-80" v-if="g.subject_name">{{ g.subject_name }}</div>
              </div>
              <div class="col-auto">
                <q-chip color="secondary" text-color="white" dense icon="people">
                  {{ g.student_count || 0 }} {{ $t('groupsPage.students') }}
                </q-chip>
              </div>
            </div>
          </q-card-section>

          <q-card-section class="q-pa-md">
            <div class="text-body2 text-grey-8 q-mb-sm" v-if="g.description">{{ g.description }}</div>
            <div class="text-caption text-grey-6" v-if="g.teacher_name">
              <q-icon name="person" class="q-mr-xs" />{{ $t('groupsPage.teacher') }}: {{ g.teacher_name }}
            </div>
            <div class="text-caption text-grey-6">
              <q-icon name="calendar_today" class="q-mr-xs" />{{ $t('groupsPage.academicYear') }}: {{ g.academic_year }}
            </div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="right">
            <q-btn flat color="info" icon="group" :label="$t('groupsPage.students')" @click="viewStudents(g)" />
            <q-btn flat color="warning" icon="edit" :label="$t('common.edit') || $t('groupsPage.editGroup')" @click="openEditDialog(g)" />
            <q-btn flat color="negative" icon="delete" :label="$t('common.delete') || $t('groupsPage.deleteGroup')" @click="deleteGroup(g.id)" />
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- Dialog Crea / Modifica Gruppo -->
    <q-dialog v-model="showDialog" persistent>
      <q-card style="width: min(500px, 95vw); max-width: 95vw;" class="q-pa-md">
        <q-form @submit="submitGroup" greedy>
          <q-card-section class="row items-center">
            <div class="text-h6">{{ isEdit ? $t('groupsPage.editGroup') : $t('groupsPage.newGroup') }}</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
          </q-card-section>

          <q-card-section class="q-gutter-md">
            <q-input
              v-model="form.name"
              :label="$t('groupsPage.groupNamePlaceholder')"
              outlined
              density="compact"
              :rules="[val => (!!val && val.trim().length > 0) || $t('common.requiredField') || 'Inserisci il nome del gruppo']"
            />
            <q-input v-model="form.description" :label="$t('groupsPage.description')" outlined type="textarea" rows="2" />
          </q-card-section>

          <q-card-actions align="right">
            <q-btn flat :label="$t('common.cancel') || 'Annulla'" v-close-popup />
            <q-btn color="primary" type="submit" :label="isEdit ? ($t('groupsPage.save') || $t('common.save')) : ($t('groupsPage.create') || $t('common.create'))" :loading="submitting" />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>

    <!-- Dialog Studenti del Gruppo -->
    <q-dialog v-model="showStudentsDialog">
      <q-card style="width: min(550px, 95vw); max-width: 95vw;" class="q-pa-md">
        <q-card-section class="row items-center">
          <div class="text-h6">{{ $t('groupsPage.studentListTitle', { name: selectedGroup?.name || '' }) }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section>
          <q-list separator v-if="selectedGroupStudents.length > 0">
            <q-item v-for="s in selectedGroupStudents" :key="s.student_id">
              <q-item-section avatar>
                <q-avatar color="primary" text-color="white" icon="person" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ s.last_name }} {{ s.first_name }}</q-item-label>
                <q-item-label caption v-if="s.class_name">{{ $t('common.class') || 'Classe' }}: {{ s.class_name }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-btn flat round color="negative" icon="person_remove" @click="removeStudent(selectedGroup.id, s.student_id)" />
              </q-item-section>
            </q-item>
          </q-list>
          <div v-else class="text-center text-grey-6 q-pa-md">{{ $t('groupsPage.noStudents') }}</div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import groupsService from '@/services/groupsService';
import { useAuthStore } from '@/stores/auth';

const $q = useQuasar();
const { t } = useI18n();
const authStore = useAuthStore();

const groups = ref([]);
const loading = ref(true);
const submitting = ref(false);

const showDialog = ref(false);
const isEdit = ref(false);
const currentGroupId = ref(null);

const form = ref({
  name: '',
  description: ''
});

const showStudentsDialog = ref(false);
const selectedGroup = ref(null);
const selectedGroupStudents = ref([]);

const fetchGroups = async () => {
  loading.value = true;
  try {
    const schoolID = authStore.user?.school_id || '';
    const res = await groupsService.getGroups({ school_id: schoolID, teacher_id: authStore.user?.id });
    groups.value = res.data || [];
  } catch (err) {
    $q.notify({ type: 'negative', message: t('groupsPage.loading') || 'Errore durante il caricamento dei gruppi' });
  } finally {
    loading.value = false;
  }
};

const openCreateDialog = () => {
  isEdit.value = false;
  currentGroupId.value = null;
  form.value = { name: '', description: '' };
  showDialog.value = true;
};

const openEditDialog = (g) => {
  isEdit.value = true;
  currentGroupId.value = g.id;
  form.value = { name: g.name, description: g.description || '' };
  showDialog.value = true;
};

const submitGroup = async () => {
  if (!form.value.name) {
    $q.notify({ type: 'warning', message: t('common.requiredField') || 'Inserisci un nome per il gruppo' });
    return;
  }
  submitting.value = true;
  try {
    if (isEdit.value) {
      await groupsService.updateGroup(currentGroupId.value, form.value);
      $q.notify({ type: 'positive', message: t('groupsPage.groupSaved') });
    } else {
      await groupsService.createGroup({
        school_id: authStore.user?.school_id || '',
        name: form.value.name,
        teacher_id: authStore.user?.id,
        description: form.value.description
      });
      $q.notify({ type: 'positive', message: t('groupsPage.groupCreated') });
    }
    showDialog.value = false;
    fetchGroups();
  } catch (err) {
    $q.notify({ type: 'negative', message: (t('common.error') || 'Errore: ') + (err.response?.data?.error || err.message) });
  } finally {
    submitting.value = false;
  }
};

const deleteGroup = (id) => {
  $q.dialog({
    title: t('groupsPage.confirmDelete'),
    message: t('groupsPage.confirmDeleteDesc'),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await groupsService.deleteGroup(id);
      $q.notify({ type: 'positive', message: t('groupsPage.groupDeleted') });
      fetchGroups();
    } catch (err) {
      $q.notify({ type: 'negative', message: t('groupsPage.deleteError') });
    }
  });
};

const viewStudents = async (g) => {
  selectedGroup.value = g;
  try {
    const res = await groupsService.getGroup(g.id);
    selectedGroupStudents.value = res.data?.students || [];
    showStudentsDialog.value = true;
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') || 'Impossibile caricare gli studenti del gruppo' });
  }
};

const removeStudent = async (groupId, studentId) => {
  try {
    await groupsService.removeStudent(groupId, studentId);
    $q.notify({ type: 'positive', message: t('common.success') || 'Studente rimosso dal gruppo' });
    viewStudents(selectedGroup.value);
    fetchGroups();
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') || 'Errore nella rimozione dello studente' });
  }
};

onMounted(() => {
  fetchGroups();
});
</script>
