<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
       <div class="row items-center q-gutter-sm">
         <div class="text-h5 text-weight-bold text-slate-800">{{ $t('classRegister.title') }}</div>
         <q-chip v-if="isSubstitutionMode" color="deep-orange" text-color="white" dense class="text-weight-bold">
           <q-icon name="swap_horiz" class="q-mr-xs" />{{ $t('classRegister.typeSubstitution') }}
         </q-chip>
       </div>
       <div class="row q-gutter-sm items-center">
           <q-btn
             :color="isSubstitutionMode ? 'deep-orange' : 'primary'"
             :icon="isSubstitutionMode ? 'swap_horiz' : 'school'"
             :label="isSubstitutionMode ? $t('classRegister.backToMyClasses') : $t('classRegister.substitutionBtn')"
             unelevated dense
             class="rounded-lg text-weight-bold"
             @click="toggleSubstitutionMode"
           />
           <q-btn color="secondary" icon="download" :label="$t('classRegister.exportCSV')" unelevated dense @click="exportCSV" />
       </div>
    </div>

    <!-- Filters Row -->
    <div class="row q-col-gutter-sm q-mb-md items-center">
      <div class="col-12 col-sm-auto">
        <q-input dense outlined v-model="date" type="date" :label="$t('classRegister.dateLabel')" bg-color="white" style="min-width: 150px" @update:model-value="fetchData" />
      </div>
      <div class="col-12 col-sm-auto">
        <q-select
           dense outlined
           v-model="selectedHour"
           :options="[1,2,3,4,5,6,7,8]"
           :label="$t('classRegister.hourLabel')"
           bg-color="white"
           style="min-width: 100px"
           @update:model-value="fetchData"
        />
      </div>
      <div class="col-12 col-sm">
        <q-select
           dense outlined
           v-model="selectedClass"
           :options="availableClassOptions"
           option-label="label"
           option-value="id"
           :label="$t('classRegister.selectClass')"
           bg-color="white"
           style="min-width: 240px"
           @update:model-value="onClassChange"
        />
      </div>
    </div>

    <!-- Substitution Banner -->
    <q-banner v-if="isSubstitutionMode" class="bg-amber-1 text-amber-10 rounded-xl border border-amber-300 q-mb-md shadow-soft">
      <template v-slot:avatar>
        <q-icon name="swap_horiz" color="amber-9" size="24px" />
      </template>
      <div class="text-weight-bold">{{ $t('classRegister.substitutionBannerTitle', { class: selectedClass?.label || $t('classRegister.selectClass') }) }}</div>
      <div class="text-caption">{{ $t('classRegister.substitutionBannerDesc') }}</div>
    </q-banner>

    <!-- Read-Only Banner: another teacher has signed this hour -->
    <q-banner v-if="isReadOnly && !isSubstitutionMode" class="bg-blue-grey-1 text-blue-grey-9 rounded-xl border border-blue-grey-3 q-mb-md shadow-soft">
      <template v-slot:avatar>
        <q-icon name="lock" color="blue-grey-7" size="24px" />
      </template>
      <div class="text-weight-bold">{{ $t('classRegister.readOnlyBannerTitle', { hour: selectedHour, teacher: signedByTeacherName }) }}</div>
      <div class="text-caption">{{ $t('classRegister.readOnlyBannerDesc') }}</div>
    </q-banner>

    <!-- Summary Cards (compact) -->
    <div class="row q-col-gutter-sm q-mb-md">
        <div class="col-6 col-sm-3">
            <q-card flat bordered class="bg-green-1">
                <q-card-section class="q-pa-sm text-center row items-center justify-center q-gutter-xs">
                    <q-icon name="check_circle" size="20px" color="green-8" />
                    <span class="text-caption text-green-9 text-weight-bold">{{ $t('classRegister.present') }}</span>
                    <span class="text-h6 text-green-8 text-weight-bolder">{{ stats.present }}</span>
                    <q-badge v-if="unjustifiedCount > 0" color="negative" class="q-ml-xs" floating>!</q-badge>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-6 col-sm-3">
             <q-card flat bordered class="bg-red-1">
                <q-card-section class="q-pa-sm text-center row items-center justify-center q-gutter-xs">
                    <q-icon name="cancel" size="20px" color="red-8" />
                    <span class="text-caption text-red-9 text-weight-bold">{{ $t('classRegister.absent') }}</span>
                    <span class="text-h6 text-red-8 text-weight-bolder">{{ stats.absent }}</span>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-6 col-sm-3">
             <q-card flat bordered class="bg-orange-1">
                <q-card-section class="q-pa-sm text-center row items-center justify-center q-gutter-xs">
                    <q-icon name="schedule" size="20px" color="orange-8" />
                    <span class="text-caption text-orange-9 text-weight-bold">{{ $t('classRegister.late') }}</span>
                    <span class="text-h6 text-orange-8 text-weight-bolder">{{ stats.late }}</span>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-6 col-sm-3">
             <q-card flat bordered class="bg-purple-1">
                <q-card-section class="q-pa-sm text-center row items-center justify-center q-gutter-xs">
                    <q-icon name="output" size="20px" color="purple-8" />
                    <span class="text-caption text-purple-9 text-weight-bold">{{ $t('classRegister.earlyExit') }}</span>
                    <span class="text-h6 text-purple-8 text-weight-bolder">{{ stats.early }}</span>
                </q-card-section>
            </q-card>
        </div>
    </div>

    <!-- Unjustified Warning Banner -->
    <q-banner v-if="unjustifiedStudents.length > 0" class="bg-red-1 text-red-9 rounded-xl border border-red-3 q-mb-md" dense>
      <template v-slot:avatar><q-icon name="warning" color="negative" size="22px" /></template>
      <div class="text-weight-bold text-caption">
        {{ $t('classRegister.unjustifiedBanner', { count: unjustifiedStudents.length }) }}
        <span v-for="(s, i) in unjustifiedStudents" :key="s.id">{{ s.last_name }} {{ s.first_name }}<span v-if="i < unjustifiedStudents.length - 1">, </span></span>
      </div>
    </q-banner>

    <!-- Section 1: Firma Lezione & Registro Lezioni svolte -->
    <q-card class="q-mb-md shadow-2 rounded-xl">
        <q-card-section class="bg-slate-800 text-white row items-center justify-between q-py-sm">
            <div class="text-subtitle2 text-weight-bold row items-center">
                <q-icon name="edit_note" class="q-mr-sm" size="20px" />
                {{ $t('classRegister.signLessonTitle', { hour: selectedHour, date: date }) }}
            </div>
            <div class="row items-center q-gutter-xs">
                <q-btn
                    v-if="canDeleteCurrentSignature"
                    unelevated
                    dense
                    color="negative"
                    text-color="white"
                    icon="delete_forever"
                    :label="$t('classRegister.deleteSignature')"
                    class="text-weight-bold q-mr-sm q-px-sm rounded-md"
                    @click="deleteUnifiedRecord"
                    :loading="saving"
                >
                    <q-tooltip>{{ $t('classRegister.deleteSignatureTooltip') }}</q-tooltip>
                </q-btn>
                <q-badge color="primary" class="text-weight-bold">
                    {{ isSubstitutionMode ? $t('classRegister.typeSubstitution') : getSubjectName(lessonSubjectId) }}
                </q-badge>
            </div>
        </q-card-section>

        <q-card-section class="q-pa-md">
            <div class="row q-col-gutter-md">
                <!-- Hour selector inside the card -->
                <div class="col-12 col-md-2">
                    <q-select
                        v-model="selectedHour"
                        :options="[1,2,3,4,5,6,7,8]"
                        :label="$t('classRegister.hourToSign')"
                        outlined
                        dense
                        @update:model-value="fetchData"
                    >
                        <template v-slot:prepend>
                            <q-icon name="schedule" />
                        </template>
                    </q-select>
                </div>
                <!-- Subject selector inside the lesson card (not in top header) -->
                <div v-if="!isSubstitutionMode" class="col-12 col-md-3">
                    <q-select
                        v-model="lessonSubjectId"
                        :options="availableSubjectOptions"
                        option-label="subject_name"
                        option-value="subject_id"
                        emit-value
                        map-options
                        :label="$t('classRegister.subjectLabel')"
                        outlined
                        dense
                    />
                </div>
                <div class="col-12" :class="isSubstitutionMode ? 'col-md-10' : 'col-md-7'">
                    <q-input
                        v-model="lessonTopic"
                        :label="$t('classRegister.topicLabel')"
                        outlined
                        dense
                        :placeholder="$t('classRegister.topicPlaceholder')"
                        :readonly="isReadOnly"
                        :bg-color="isReadOnly ? 'grey-2' : 'white'"
                    />
                </div>
                <div class="col-12 col-md-3">
                    <q-select
                        v-model="lessonType"
                        :options="lessonTypeOptions"
                        :label="$t('classRegister.lessonTypeLabel')"
                        outlined
                        dense
                        :readonly="isReadOnly"
                    />
                </div>
                <div class="col-12 col-md-4">
                    <q-select
                        v-model="activityType"
                        :options="activityTypeOptions"
                        option-value="value"
                        option-label="label"
                        emit-value
                        map-options
                        :label="$t('classRegister.activityTypeLabel')"
                        outlined
                        dense
                        :readonly="isReadOnly"
                    >
                        <template v-slot:option="scope">
                            <q-item v-bind="scope.itemProps">
                                <q-item-section avatar>
                                    <q-icon :name="scope.opt.icon" :color="scope.opt.color" />
                                </q-item-section>
                                <q-item-section>
                                    <q-item-label>{{ scope.opt.label }}</q-item-label>
                                    <q-item-label caption>{{ scope.opt.caption }}</q-item-label>
                                </q-item-section>
                            </q-item>
                        </template>
                    </q-select>
                </div>
                <div class="col-12 col-md-2 row items-center">
                    <q-toggle v-model="isCoTeaching" :label="$t('classRegister.coTeaching')" color="deep-purple" dense :disable="isReadOnly" />
                </div>
                <div class="col-12 col-md-3">
                    <q-input v-model="lessonNotes" :label="$t('classRegister.notesLabel')" outlined dense autogrow :placeholder="$t('classRegister.notesPlaceholder')" :readonly="isReadOnly" />
                </div>
                <div v-if="isPctoOrOrientamento" class="col-12">
                    <q-banner class="bg-blue-1 text-blue-9 rounded-borders" dense>
                        <template v-slot:avatar><q-icon name="info" color="blue-7" /></template>
                        {{ $t('classRegister.noGradesWarning', { type: getActivityTypeLabel(activityType) }) }}
                    </q-banner>
                </div>
            </div>

            <!-- Homework Accordion -->
            <div class="q-mt-sm">
                <q-toggle v-model="assignHomework" :label="$t('classRegister.assignHomework')" color="orange" dense :disable="isReadOnly" />
                <div v-if="assignHomework" class="row q-col-gutter-md q-mt-xs">
                    <div class="col-12 col-md-8">
                        <q-input v-model="homeworkDesc" :label="$t('classRegister.homeworkDesc')" outlined dense autogrow :placeholder="$t('classRegister.homeworkDescPlaceholder')" />
                    </div>
                    <div class="col-12 col-md-4">
                        <q-input v-model="homeworkDue" type="date" :label="$t('classRegister.homeworkDueDate')" outlined dense />
                    </div>
                </div>
            </div>
        </q-card-section>

        <!-- Daily Lessons Timeline -->
        <q-separator />
        <div class="q-pa-sm bg-slate-50">
            <div class="text-caption text-weight-bold text-slate-700 q-mb-xs q-px-sm">{{ $t('classRegister.todayLessonsTitle') }}</div>
            <div v-if="dailyLessons.length > 0" class="row q-gutter-xs q-px-sm">
                <div v-for="l in dailyLessons" :key="l.id" class="col-auto">
                    <q-chip dense outline :color="String(l.hour) === String(selectedHour) ? 'indigo-9' : 'indigo-5'" icon="event_note" class="bg-white">
                        {{ $t('classRegister.hourLabel') }} {{ l.hour }}: {{ l.topic || '—' }}
                        <q-badge v-if="l.activity_type && l.activity_type !== 'standard'" :color="getActivityTypeColor(l.activity_type)" class="q-ml-xs text-caption">
                            <q-icon :name="getActivityTypeIcon(l.activity_type)" size="12px" class="q-mr-xs" />{{ getActivityTypeLabel(l.activity_type) }}
                        </q-badge>
                        <q-tooltip>
                            Docente: {{ l.teacher_name || 'Docente' }}<br>
                            Materia: {{ getSubjectName(l.subject_id) }}<br>
                            Tipo: {{ l.type }}<br>
                            Attività: {{ getActivityTypeLabel(l.activity_type || 'standard') }}
                        </q-tooltip>
                    </q-chip>
                </div>
            </div>
            <div v-else class="q-px-sm text-caption text-grey-6">{{ $t('classRegister.noLessonsToday') }}</div>
        </div>
    </q-card>

    <!-- Section 2: Attendance Table with Hourly Timeline Column -->
    <q-card class="shadow-2 rounded-xl">
        <q-toolbar class="bg-grey-2 text-grey-9">
            <q-toolbar-title class="text-subtitle2 row items-center">
                <q-icon name="how_to_reg" class="q-mr-xs" color="primary" />
                <span>{{ $t('classRegister.attendanceSectionTitle') }} — {{ $t('classRegister.hourLabel') }} {{ selectedHour }}</span>
                <q-chip dense :color="markedCountChipColor" text-color="white" class="q-ml-sm">
                    {{ markedCount }}/{{ students.length }}
                </q-chip>
                <q-chip v-if="lastAutosaveTime" dense color="grey-7" text-color="white" icon="cloud_done" class="q-ml-sm text-caption">
                    Bozza {{ lastAutosaveTime }}
                </q-chip>
            </q-toolbar-title>
            <q-btn flat dense icon="check_circle" :label="$t('classRegister.markAllPresent')" color="primary" @click="markAllPresent" :disable="loading || isReadOnly" />
        </q-toolbar>

        <div v-if="loading" class="q-pa-md">
            <SkeletonTable :rows="8" :cols="4" />
        </div>

        <q-list separator v-else>
            <q-item v-for="student in students" :key="student.id" class="q-py-sm transition-bg" :class="getRowClass(student.status)">
                <q-item-section avatar style="min-width: 42px">
                    <q-avatar
                        size="38px"
                        color="indigo-1"
                        text-color="indigo-9"
                        class="text-weight-bold shadow-soft cursor-pointer"
                        @click="openStudentPanel(student)"
                    >
                        {{ student.first_name ? student.first_name.charAt(0) : '?' }}
                        <q-tooltip>{{ $t('classRegister.studentInfo') }}</q-tooltip>
                    </q-avatar>
                </q-item-section>

                <!-- Name + badges column -->
                <q-item-section style="min-width: 160px">
                    <q-item-label class="text-weight-bold text-body2">
                        {{ student.last_name }} {{ student.first_name }}
                        <q-icon
                          v-if="student.hasUnjustified"
                          name="warning"
                          color="negative"
                          size="16px"
                          class="q-ml-xs"
                        >
                          <q-tooltip>{{ $t('classRegister.unjustifiedBanner', { count: 1 }) }}</q-tooltip>
                        </q-icon>
                    </q-item-label>
                    <!-- Hourly presence indicator (mini timeline) -->
                    <div class="row items-center q-gutter-xs q-mt-xs">
                        <q-badge
                          v-for="h in 8"
                          :key="h"
                          :color="getHourBadgeColor(student, h)"
                          :label="String(h)"
                          class="text-weight-bold"
                          style="min-width: 18px; font-size: 10px; padding: 2px 4px"
                        >
                          <q-tooltip>{{ getHourLabel(student, h) }}</q-tooltip>
                        </q-badge>
                    </div>
                    <!-- Status / time badges -->
                    <div class="row items-center q-gutter-xs q-mt-xs">
                        <q-badge v-if="student.status === 'Present'" color="positive" class="q-px-sm q-py-xs">
                            <q-icon name="check_circle" class="q-mr-xs" size="12px" /> {{ $t('classRegister.present') }}
                        </q-badge>
                        <q-badge v-else-if="student.status === 'OutOfClass'" color="teal" class="q-px-sm q-py-xs">
                            <q-icon name="meeting_room" class="q-mr-xs" size="12px" /> Fuori Aula
                        </q-badge>
                        <q-badge v-else-if="student.status === 'Absent'" color="negative" class="q-px-sm q-py-xs">
                            <q-icon name="cancel" class="q-mr-xs" size="12px" /> {{ $t('classRegister.absent') }}
                        </q-badge>
                        <q-badge v-else-if="student.status === 'Late'" color="warning" text-color="black" class="q-px-sm q-py-xs text-weight-bold">
                            <q-icon name="schedule" class="q-mr-xs" size="12px" />
                            {{ student.entry_time ? `${$t('classRegister.tableHeaderEntryTime')} ${student.entry_time}` : $t('classRegister.late') }}
                        </q-badge>
                        <q-badge v-else-if="student.status === 'LeftEarly'" color="purple" class="q-px-sm q-py-xs text-weight-bold">
                            <q-icon name="output" class="q-mr-xs" size="12px" />
                            {{ student.exit_time ? `${$t('classRegister.tableHeaderExitTime')} ${student.exit_time}` : $t('classRegister.earlyExit') }}
                        </q-badge>
                    </div>
                </q-item-section>

                <!-- Status toggle -->
                <q-item-section>
                    <q-btn-toggle
                        v-model="student.status"
                        flat dense no-caps
                        :disable="isReadOnly"
                        :options="[
                            {icon: 'check', value: 'Present', slot: 'present'},
                            {icon: 'meeting_room', value: 'OutOfClass', slot: 'outofclass'},
                            {icon: 'close', value: 'Absent', slot: 'absent'},
                            {icon: 'schedule', value: 'Late', slot: 'late'},
                            {icon: 'logout', value: 'LeftEarly', slot: 'early'}
                        ]"
                    >
                        <template v-slot:present><q-tooltip>{{ $t('classRegister.present') }}</q-tooltip></template>
                        <template v-slot:outofclass><q-tooltip>Fuori Aula</q-tooltip></template>
                        <template v-slot:absent><q-tooltip>{{ $t('classRegister.absent') }}</q-tooltip></template>
                        <template v-slot:late><q-tooltip>{{ $t('classRegister.late') }}</q-tooltip></template>
                        <template v-slot:early><q-tooltip>{{ $t('classRegister.earlyExit') }}</q-tooltip></template>
                    </q-btn-toggle>
                </q-item-section>

                <!-- Late Entry Time -->
                <q-item-section v-if="student.status === 'Late'" side style="min-width: 130px">
                     <q-input
                        v-model="student.entry_time"
                        type="time"
                        dense outlined
                        :label="$t('classRegister.tableHeaderEntryTime')"
                        bg-color="white"
                     />
                </q-item-section>

                <!-- Early Exit Time -->
                <q-item-section v-if="student.status === 'LeftEarly'" side style="min-width: 130px">
                     <q-input
                        v-model="student.exit_time"
                        type="time"
                        dense outlined
                        :label="$t('classRegister.tableHeaderExitTime')"
                        bg-color="white"
                     />
                </q-item-section>

                <q-item-section side>
                    <q-btn round flat icon="note_add" color="grey-7" @click="openNoteDialog(student)">
                        <q-tooltip>{{ $t('classRegister.addDisciplinaryNote') }}</q-tooltip>
                    </q-btn>
                </q-item-section>
            </q-item>

            <q-item v-if="students.length === 0" class="text-center text-grey">
                <q-item-section>{{ $t('classRegister.noStudentsFound') }}</q-item-section>
            </q-item>
        </q-list>

        <q-card-actions align="between" class="bg-grey-1 q-pa-md">
            <div>
                <q-btn
                    v-if="canDeleteCurrentSignature"
                    :label="$t('classRegister.deleteSignatureAndAttendance')"
                    color="negative"
                    unelevated
                    icon="delete_outline"
                    class="rounded-lg text-weight-bold"
                    @click="deleteUnifiedRecord"
                    :loading="saving"
                />
            </div>
            <div>
                <q-btn
                  :label="isReadOnly ? $t('classRegister.readOnlySaveBtn') : (currentHourLesson ? $t('classRegister.updateSaveBtn') : $t('classRegister.saveBtn'))"
                  :color="isReadOnly ? 'grey-6' : isSubstitutionMode ? 'deep-orange' : 'primary'"
                  size="md"
                  :icon="isReadOnly ? 'lock' : 'cloud_done'"
                  class="rounded-lg text-weight-bold shadow-soft"
                  @click="saveUnifiedRecord"
                  :loading="saving"
                  :disable="!selectedClass || isReadOnly"
                />
            </div>
        </q-card-actions>
    </q-card>

    <!-- Note Dialog -->
    <NoteDialog
        v-if="selectedClass"
        v-model="showNoteDialog"
        :student="selectedStudentForNote"
        :class-id="String(typeof selectedClass === 'object' ? selectedClass.id : selectedClass)"
    />

    <!-- Student Detail Dialog -->
    <q-dialog v-model="showStudentPanel" maximized-if-mobile>
        <q-card style="min-width: 340px; max-width: 520px; width: 100%">
            <q-bar class="bg-indigo-8 text-white">
                <q-icon name="person" class="q-mr-sm" />
                <span class="text-subtitle2 text-weight-bold">
                    {{ panelStudent?.last_name }} {{ panelStudent?.first_name }}
                </span>
                <q-space />
                <q-btn flat round dense icon="close" v-close-popup />
            </q-bar>

            <q-card-section class="q-pa-md">
                <!-- Loading -->
                <div v-if="panelLoading" class="row justify-center q-pa-lg">
                    <q-spinner color="primary" size="40px" />
                </div>

                <template v-else>
                    <!-- Personal Info -->
                    <div class="text-caption text-weight-bold text-grey-6 q-mb-xs text-uppercase letter-spacing-wide">Dati Anagrafici</div>
                    <q-list bordered separator rounded class="q-mb-md">
                        <q-item dense>
                            <q-item-section avatar><q-icon name="badge" color="indigo" /></q-item-section>
                            <q-item-section>
                                <q-item-label caption>Nome completo</q-item-label>
                                <q-item-label>{{ panelStudentInfo?.last_name }} {{ panelStudentInfo?.first_name }}</q-item-label>
                            </q-item-section>
                        </q-item>
                        <q-item dense v-if="panelStudentInfo?.fiscal_code">
                            <q-item-section avatar><q-icon name="fingerprint" color="indigo" /></q-item-section>
                            <q-item-section>
                                <q-item-label caption>Codice Fiscale</q-item-label>
                                <q-item-label class="text-mono">{{ panelStudentInfo.fiscal_code }}</q-item-label>
                            </q-item-section>
                        </q-item>
                        <q-item dense v-if="panelStudentInfo?.class_name">
                            <q-item-section avatar><q-icon name="class" color="indigo" /></q-item-section>
                            <q-item-section>
                                <q-item-label caption>Classe</q-item-label>
                                <q-item-label>{{ panelStudentInfo.class_name }}</q-item-label>
                            </q-item-section>
                        </q-item>
                        <q-item dense v-if="panelStudentInfo?.email">
                            <q-item-section avatar><q-icon name="email" color="indigo" /></q-item-section>
                            <q-item-section>
                                <q-item-label caption>Email</q-item-label>
                                <q-item-label>{{ panelStudentInfo.email }}</q-item-label>
                            </q-item-section>
                        </q-item>
                        <q-item dense v-if="panelStudentInfo?.phone_number">
                            <q-item-section avatar><q-icon name="phone" color="indigo" /></q-item-section>
                            <q-item-section>
                                <q-item-label caption>Telefono</q-item-label>
                                <q-item-label>{{ panelStudentInfo.phone_number }}</q-item-label>
                            </q-item-section>
                        </q-item>
                        <q-item dense v-if="panelStudentInfo?.date_of_birth">
                            <q-item-section avatar><q-icon name="cake" color="indigo" /></q-item-section>
                            <q-item-section>
                                <q-item-label caption>Data di Nascita</q-item-label>
                                <q-item-label>{{ formatDate(panelStudentInfo.date_of_birth) }}</q-item-label>
                            </q-item-section>
                        </q-item>
                    </q-list>

                    <!-- Attendance Summary -->
                    <div class="text-caption text-weight-bold text-grey-6 q-mb-xs text-uppercase">Riepilogo Presenze</div>
                    <div class="row q-col-gutter-sm q-mb-md">
                        <div class="col-6">
                            <q-card flat bordered class="text-center q-pa-sm">
                                <div class="text-h5 text-negative text-weight-bold">{{ panelSummary?.total_absences ?? '—' }}</div>
                                <div class="text-caption text-grey-7">Assenze Totali</div>
                            </q-card>
                        </div>
                        <div class="col-6">
                            <q-card flat bordered class="text-center q-pa-sm">
                                <div class="text-h5 text-warning text-weight-bold">{{ panelSummary?.total_lates ?? '—' }}</div>
                                <div class="text-caption text-grey-7">Ritardi</div>
                            </q-card>
                        </div>
                        <div class="col-6">
                            <q-card flat bordered class="text-center q-pa-sm">
                                <div class="text-h5 text-purple text-weight-bold">{{ panelSummary?.total_early_exits ?? '—' }}</div>
                                <div class="text-caption text-grey-7">Uscite Anticipate</div>
                            </q-card>
                        </div>
                        <div class="col-6">
                            <q-card flat bordered class="text-center q-pa-sm">
                                <div class="text-h5 text-positive text-weight-bold">{{ panelSummary?.justified_count ?? '—' }}</div>
                                <div class="text-caption text-grey-7">Giustificate</div>
                            </q-card>
                        </div>
                    </div>

                    <!-- Absence Rate + Risk -->
                    <div v-if="panelSummary" class="q-mb-sm">
                        <div class="row items-center justify-between q-mb-xs">
                            <span class="text-caption text-grey-7">Tasso assenza</span>
                            <span class="text-caption text-weight-bold" :class="panelSummary.absence_rate > 25 ? 'text-negative' : panelSummary.absence_rate > 10 ? 'text-warning' : 'text-positive'">
                                {{ panelSummary.absence_rate?.toFixed(1) }}%
                            </span>
                        </div>
                        <q-linear-progress
                            :value="(panelSummary.absence_rate || 0) / 100"
                            :color="panelSummary.absence_rate > 25 ? 'negative' : panelSummary.absence_rate > 10 ? 'warning' : 'positive'"
                            rounded
                            size="8px"
                            class="q-mb-xs"
                        />
                        <q-chip
                            dense
                            :color="panelSummary.risk_level === 'high' ? 'negative' : panelSummary.risk_level === 'medium' ? 'warning' : 'positive'"
                            text-color="white"
                            :icon="panelSummary.risk_level === 'high' ? 'warning' : panelSummary.risk_level === 'medium' ? 'info' : 'check_circle'"
                        >
                            Rischio: {{ panelSummary.risk_level === 'high' ? 'ALTO' : panelSummary.risk_level === 'medium' ? 'MEDIO' : 'BASSO' }}
                        </q-chip>
                    </div>

                    <!-- Today's attendance for this student -->
                    <div class="text-caption text-weight-bold text-grey-6 q-mt-md q-mb-xs text-uppercase">Presenze Oggi per Ora</div>
                    <div class="row q-gutter-xs">
                        <q-badge
                            v-for="h in 8" :key="h"
                            :color="getHourBadgeColor(panelStudent, h)"
                            :label="String(h) + 'ª'"
                            class="text-weight-bold"
                            style="font-size: 11px; padding: 4px 8px"
                        >
                            <q-tooltip>{{ getHourLabel(panelStudent, h) }}</q-tooltip>
                        </q-badge>
                    </div>
                </template>
            </q-card-section>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useClassesStore } from '@/stores/classes'
import { useGradesStore } from '@/stores/grades'
import { useAuthStore } from '@/stores/auth'
import { attendanceService } from 'src/services/attendanceService'
import { lessonService } from 'src/services/lessonService'
import api from '@/services/api'
import NoteDialog from 'src/components/Teacher/NoteDialog.vue'
import SkeletonTable from '@/components/Common/SkeletonTable.vue'
import { useSchoolYearStore } from '@/stores/schoolYear'

const $q = useQuasar()
let t = (key, fallback) => (typeof fallback === 'string' ? fallback : key)
try {
  const i18nInstance = useI18n()
  if (i18nInstance && i18nInstance.t) {
    t = i18nInstance.t
  }
} catch (e) {
  // Fallback for isolated unit tests without vue-i18n app plugin
}
const classesStore = useClassesStore()
const gradesStore = useGradesStore()
const authStore = useAuthStore()
const schoolYearStore = useSchoolYearStore()

// Current teacher ID from auth token
const currentTeacherId = computed(() => authStore.user?.id || null)

const date = ref(new Date().toISOString().split('T')[0])
const selectedClass = ref(null)
const students = ref([])
const isSubstitutionMode = ref(false)
const allSchoolClasses = ref([])

// Lesson Form State
const lessonTopic = ref('')
const lessonType = ref('Frontale')
const activityType = ref('standard')
const lessonSubjectId = ref(null)
const isCoTeaching = ref(false)
const lessonNotes = ref('')
const assignHomework = ref(false)
const homeworkDesc = ref('')
const homeworkDue = ref('')

// All attendance records for this class today (all hours) for timeline
const allTodayAttendance = ref([])

const lessonTypeOptions = computed(() => [
  t('classRegister.typeCurricular'),
  t('classRegister.typeSubstitution'),
  t('classRegister.activityLab'),
  t('classRegister.typeExtracurricular'),
  t('classRegister.typeCoTeaching'),
  t('classRegister.activityStandard')
])

/** Opzioni per la tipologia di attività */
const activityTypeOptions = computed(() => [
  { value: 'standard',          label: t('classRegister.activityStandard'),          icon: 'menu_book',   color: 'primary',     caption: t('classRegister.activityStandardCap') },
  { value: 'substitution',      label: t('classRegister.typeSubstitution'),          icon: 'swap_horiz',   color: 'deep-orange', caption: t('classRegister.substitutionBannerDesc') },
  { value: 'pcto',              label: t('classRegister.activityPcto'),              icon: 'work',         color: 'deep-purple', caption: t('classRegister.activityPctoCap') },
  { value: 'orientamento',      label: t('classRegister.activityOrientamento'),      icon: 'explore',      color: 'teal',        caption: t('classRegister.activityOrientamentoCap') },
  { value: 'pcto_orientamento', label: `${t('classRegister.activityPcto')} - ${t('classRegister.activityOrientamento')}`, icon: 'hub', color: 'indigo-8', caption: t('classRegister.activityPctoCap') },
  { value: 'lab',               label: t('classRegister.activityLab'),               icon: 'biotech',      color: 'green',       caption: t('classRegister.activityLabCap') },
  { value: 'recupero',          label: t('classRegister.activityRecupero'),          icon: 'healing',      color: 'purple',      caption: t('classRegister.activityRecuperoCap') }
])

const getActivityTypeColor = (type) => {
  const opt = activityTypeOptions.value?.find(o => o.value === type)
  return opt ? opt.color : 'grey'
}
const getActivityTypeIcon = (type) => {
  const opt = activityTypeOptions.value?.find(o => o.value === type)
  return opt ? opt.icon : 'menu_book'
}
const getActivityTypeLabel = (type) => {
  const opt = activityTypeOptions.value?.find(o => o.value === type)
  return opt ? opt.label : (type || 'Standard')
}
const isPctoOrOrientamento = computed(() =>
  activityType.value === 'pcto' ||
  activityType.value === 'orientamento' ||
  activityType.value === 'pcto_orientamento'
)

const availableClassOptions = computed(() => {
  const source = isSubstitutionMode.value ? allSchoolClasses.value : classesStore.classes
  return source.map(c => {
    let nameText = c.name || `Classe ${c.id}`
    if (c.section && !nameText.endsWith(c.section)) nameText += c.section
    if (c.articolazione) nameText += ` - ${c.articolazione}`
    return { ...c, label: c.label || nameText }
  })
})

const isCivicaSubject = (s) => {
  const name = (s.subject_name || s.name || '').toLowerCase()
  return name.includes('civica') || name.includes('educazione civica') || name.includes('ed. civica')
}

const isAssignedToCurrentTeacher = (s, user) => {
  if (!user) return false
  const currentUserId = String(user.id || '')
  const teacherId = user.teacher_id ? String(user.teacher_id) : ''
  const sTeacherId = s.teacher_id ? String(s.teacher_id) : ''
  const sTeacherUserId = s.teacher_user_id ? String(s.teacher_user_id) : ''

  if (sTeacherId && (sTeacherId === currentUserId || (teacherId && sTeacherId === teacherId))) {
    return true
  }
  if (sTeacherUserId && (sTeacherUserId === currentUserId || (teacherId && sTeacherUserId === teacherId))) {
    return true
  }
  if (s.teacher_name && user.last_name) {
    const tName = s.teacher_name.toLowerCase()
    const uLast = user.last_name.toLowerCase()
    const uFirst = (user.first_name || '').toLowerCase()
    if (tName.includes(uLast) && (!uFirst || tName.includes(uFirst))) {
      return true
    }
  }
  return false
}

const availableSubjectOptions = computed(() => {
  if (isSubstitutionMode.value) {
    return [{ subject_name: 'Supplenza / Compresenza', subject_id: 'supplenza' }]
  }
  const allSubjects = gradesStore.subjects || []
  const user = authStore.user
  if (!user || ['admin', 'superadmin', 'secretary'].includes(user.role)) {
    return allSubjects
  }

  // Mostra ESCLUSIVAMENTE le materie assegnate dalla segreteria al docente loggato + Educazione Civica
  return allSubjects.filter(s => {
    return isAssignedToCurrentTeacher(s, user) || isCivicaSubject(s)
  })
})

const loading = ref(false)
const saving = ref(false)
const lastAutosaveTime = ref('')

// Note Dialog State
const showNoteDialog = ref(false)
const selectedStudentForNote = ref(null)

// ── Student Detail Panel ──────────────────────────────────────
const showStudentPanel = ref(false)
const panelStudent = ref(null)
const panelStudentInfo = ref(null)
const panelSummary = ref(null)
const panelLoading = ref(false)

async function openStudentPanel(student) {
    panelStudent.value = student
    panelStudentInfo.value = null
    panelSummary.value = null
    panelLoading.value = true
    showStudentPanel.value = true
    try {
        const [infoRes, summaryRes] = await Promise.allSettled([
            api.get(`/users/${student.id}`),
            attendanceService.getStudentSummary(student.id)
        ])
        if (infoRes.status === 'fulfilled') panelStudentInfo.value = infoRes.value.data
        if (summaryRes.status === 'fulfilled') panelSummary.value = summaryRes.value
    } catch (e) {
        // individual errors handled above via allSettled
    } finally {
        panelLoading.value = false
    }
}
// ─────────────────────────────────────────────────────────────

// Formatta una data ISO (YYYY-MM-DD) in formato italiano (dd/mm/yyyy)
function formatDate(isoDate) {
    if (!isoDate) return ''
    const [y, m, d] = isoDate.split('-')
    return `${d}/${m}/${y}`
}

const selectedHour = ref(1)

const dailyLessons = ref([])

// The lesson signed for the currently selected hour (if any)
const currentHourLesson = computed(() =>
    dailyLessons.value.find(l => String(l.hour) === String(selectedHour.value)) || null
)

// Read-only if another teacher has signed this hour (and it's not substitution mode)
const isReadOnly = computed(() => {
    if (isSubstitutionMode.value) return false
    if (!currentHourLesson.value) return false
    return currentHourLesson.value.teacher_id !== currentTeacherId.value
})

const signedByTeacherName = computed(() => currentHourLesson.value?.teacher_name || 'altro docente')

const hasCurrentHourAttendance = computed(() => {
    return allTodayAttendance.value.some(r => String(r.hour) === String(selectedHour.value))
})

const canDeleteCurrentSignature = computed(() => {
    if (!selectedClass.value || isReadOnly.value) return false
    return !!currentHourLesson.value || hasCurrentHourAttendance.value
})

const unjustifiedStudents = computed(() =>
  students.value.filter(s => s.hasUnjustified)
)
const unjustifiedCount = computed(() => unjustifiedStudents.value.length)

const stats = computed(() => ({
    present: students.value.filter(s => s.status === 'Present' || s.status === 'OutOfClass').length,
    absent: students.value.filter(s => s.status === 'Absent').length,
    late: students.value.filter(s => s.status === 'Late').length,
    early: students.value.filter(s => s.status === 'LeftEarly').length,
}))

const markedCount = computed(() => students.value.filter(s => s.status).length)
const markedCountChipColor = computed(() => {
    if (markedCount.value === 0) return 'negative'
    if (markedCount.value < students.value.length) return 'warning'
    return 'positive'
})

// Hourly timeline helpers
const getHourBadgeColor = (student, hour) => {
    const rec = allTodayAttendance.value.find(r =>
        r.student_id === student.id && String(r.hour) === String(hour)
    )
    if (!rec) return 'grey-3'
    switch (rec.status) {
        case 'Present': return 'positive'
        case 'OutOfClass': return 'teal'
        case 'Absent': return 'negative'
        case 'Late': return 'warning'
        case 'LeftEarly': return 'purple'
        default: return 'grey-3'
    }
}

const getHourLabel = (student, hour) => {
    const rec = allTodayAttendance.value.find(r =>
        r.student_id === student.id && String(r.hour) === String(hour)
    )
    if (!rec) return `Ora ${hour}: non registrata`
    const statusLabels = {
        Present: 'Presente', OutOfClass: 'Fuori Aula', Absent: 'Assente',
        Late: `Ritardo${rec.entry_time ? ' ore ' + rec.entry_time : ''}`,
        LeftEarly: `Uscita anticipata${rec.exit_time ? ' ore ' + rec.exit_time : ''}`
    }
    return `Ora ${hour}: ${statusLabels[rec.status] || rec.status}`
}

const toggleSubstitutionMode = async () => {
  isSubstitutionMode.value = !isSubstitutionMode.value
  if (isSubstitutionMode.value) {
    allSchoolClasses.value = await classesStore.fetchAllSchoolClassesAndGroups()
    if (allSchoolClasses.value.length > 0) selectedClass.value = allSchoolClasses.value[0]
    lessonSubjectId.value = null
    lessonType.value = 'Supplenza'
    activityType.value = 'substitution'
    lessonTopic.value = ''
    $q.notify({ type: 'info', message: 'Modalità Supplenza attivata', timeout: 2500 })
  } else {
    await classesStore.fetchAssignedClasses()
    if (classesStore.classes.length > 0) selectedClass.value = classesStore.classes[0]
    lessonType.value = 'Frontale'
    activityType.value = 'standard'
    lessonTopic.value = ''
    await onClassChange()
  }
  fetchData()
}

const getSubjectName = (subjectId) => {
    if (!subjectId || subjectId === 'supplenza') return 'Supplenza'
    const found = gradesStore.subjects.find(s => String(s.subject_id) === String(subjectId) || String(s.id) === String(subjectId))
    return found ? found.subject_name : subjectId
}

const getRowClass = (status) => {
    switch (status) {
        case 'Present': return 'bg-green-1'
        case 'OutOfClass': return 'bg-teal-1'
        case 'Absent': return 'bg-red-1'
        case 'Late': return 'bg-orange-1'
        case 'LeftEarly': return 'bg-purple-1'
        default: return ''
    }
}

let autosaveInterval = null
const saveDraftToStorage = () => {
    if (!selectedClass.value || students.value.length === 0) return
    const classId = typeof selectedClass.value === 'object' ? selectedClass.value?.id : selectedClass.value
    const key = `attendance_draft_${classId}_${date.value}_${selectedHour.value}`
    try {
        localStorage.setItem(key, JSON.stringify({
            date: date.value, hour: selectedHour.value,
            statuses: students.value.map(s => ({ id: s.id, status: s.status, entry_time: s.entry_time, exit_time: s.exit_time })),
            timestamp: new Date().toISOString()
        }))
        lastAutosaveTime.value = new Date().toLocaleTimeString('it-IT')
    } catch (e) { /* ignore */ }
}

onMounted(async () => {
    await classesStore.fetchAssignedClasses(schoolYearStore.selectedSchoolYear)
    if (classesStore.classes.length > 0) selectedClass.value = classesStore.classes[0]
    await onClassChange()
    autosaveInterval = setInterval(saveDraftToStorage, 60000)
})

watch(() => schoolYearStore.selectedSchoolYear, async (newSY) => {
    await classesStore.fetchAssignedClasses(newSY)
    if (classesStore.classes.length > 0) {
        selectedClass.value = classesStore.classes[0]
    } else {
        selectedClass.value = null
    }
    await onClassChange()
})

onUnmounted(() => {
    if (autosaveInterval) clearInterval(autosaveInterval)
})

const onClassChange = async () => {
    if (selectedClass.value && !isSubstitutionMode.value) {
        const classId = typeof selectedClass.value === 'object' ? selectedClass.value.id : selectedClass.value
        await gradesStore.fetchClassSubjects(classId)
        if (availableSubjectOptions.value.length > 0) {
            const isSelectedAvailable = availableSubjectOptions.value.some(s => s.subject_id === lessonSubjectId.value)
            if (!lessonSubjectId.value || !isSelectedAvailable) {
                lessonSubjectId.value = availableSubjectOptions.value[0].subject_id
            }
        } else {
            lessonSubjectId.value = null
        }
    } else if (isSubstitutionMode.value) {
        lessonSubjectId.value = null
    }
    fetchData()
}

const fetchData = async () => {
    if (!selectedClass.value) return

    loading.value = true
    try {
        const classId = typeof selectedClass.value === 'object' ? selectedClass.value.id : selectedClass.value
        // 1. Fetch Students
        const usersRes = await api.get('/users', { params: { class_id: classId, role: 'student', page_size: 500 } })
        const studentList = usersRes.data.users || usersRes.data || []

        // 2. Fetch ALL attendance for this class today (all hours for timeline)
        let attendanceAll = []
        try {
            const attRes = await attendanceService.getByClass(classId, date.value)
            attendanceAll = Array.isArray(attRes.data) ? attRes.data : (attRes.data?.records || [])
        } catch (e) { /* none yet */ }
        allTodayAttendance.value = attendanceAll

        // 3. Fetch unjustified per student
        const unjustifiedSet = new Set()
        try {
            const ujRes = await api.get('/attendance/pending-justifications', { params: { class_id: classId } })
            const pending = ujRes.data || []
            pending.forEach(j => unjustifiedSet.add(j.student_id))
        } catch (e) { /* ignore */ }

        // 4. Build attendance map for current hour
        const attendanceMap = {}
        attendanceAll.forEach(r => {
            if (String(r.hour) === String(selectedHour.value)) {
                attendanceMap[r.student_id] = r
            }
        })

        // Merge students
        students.value = studentList.map(s => {
            const existing = attendanceMap[s.id]
            return {
                id: s.id,
                first_name: s.first_name,
                last_name: s.last_name,
                status: existing ? existing.status : null,
                entry_time: existing?.entry_time || '',
                exit_time: existing?.exit_time || '',
                hasUnjustified: unjustifiedSet.has(s.id),
            }
        })

        // 5. Fetch Daily Lessons for Timeline
        try {
            const lessonRes = await api.get(`/lessons/class/${classId}`, { params: { date: date.value } })
            dailyLessons.value = lessonRes.data || []
        } catch (e) { dailyLessons.value = [] }

        // Pre-fill lesson topic from existing lesson for this hour
        // Only populate editable fields if the lesson belongs to the current teacher
        const existingLesson = dailyLessons.value.find(l => String(l.hour) === String(selectedHour.value))
        if (existingLesson && existingLesson.teacher_id === currentTeacherId.value) {
            // Own lesson: pre-fill all fields
            lessonTopic.value = existingLesson.topic || ''
            lessonType.value = existingLesson.type || 'Frontale'
            activityType.value = existingLesson.activity_type || 'standard'
            isCoTeaching.value = !!existingLesson.is_co_teaching
            lessonNotes.value = existingLesson.notes || ''
            if (existingLesson.subject_id) lessonSubjectId.value = existingLesson.subject_id
        } else if (existingLesson) {
            // Another teacher's lesson: show read-only but clear form fields for display
            lessonTopic.value = existingLesson.topic || ''
            lessonType.value = existingLesson.type || 'Frontale'
            activityType.value = existingLesson.activity_type || 'standard'
            isCoTeaching.value = false
            lessonNotes.value = ''
        } else {
            // No lesson for this hour: clear topic, keep other defaults
            lessonTopic.value = ''
            lessonNotes.value = ''
            if (!isSubstitutionMode.value) {
                lessonType.value = 'Frontale'
                activityType.value = 'standard'
                isCoTeaching.value = false
            }
        }

        // Pre-fill student statuses from the PREVIOUS hour if current hour has no records
        const hasCurrentHourRecords = Object.keys(attendanceMap).length > 0
        if (!hasCurrentHourRecords && selectedHour.value > 1) {
            const prevHour = selectedHour.value - 1
            const prevHourMap = {}
            attendanceAll.forEach(r => {
                if (String(r.hour) === String(prevHour)) {
                    prevHourMap[r.student_id] = r
                }
            })
            if (Object.keys(prevHourMap).length > 0) {
                students.value = students.value.map(s => {
                    const prev = prevHourMap[s.id]
                    if (!prev) return s
                    // Present/OutOfClass/Late → Present (they arrived)
                    // Absent/LeftEarly → Absent (they're not in class)
                    let inheritedStatus
                    if (prev.status === 'Absent' || prev.status === 'LeftEarly') {
                        inheritedStatus = 'Absent'
                    } else {
                        inheritedStatus = 'Present'
                    }
                    return { ...s, status: inheritedStatus, entry_time: '', exit_time: '' }
                })
            }
        }

    } catch (error) {
        console.error(error)
    } finally {
        loading.value = false
    }
}

const markAllPresent = () => {
    $q.dialog({
        title: 'Conferma',
        message: 'Segnare tutti gli studenti come PRESENTI per questa ora?',
        cancel: true,
        persistent: true
    }).onOk(() => {
        students.value.forEach(s => { s.status = 'Present'; s.entry_time = ''; s.exit_time = '' })
    })
}

const saveUnifiedRecord = async () => {
    const unmarked = students.value.filter(s => !s.status)
    if (unmarked.length > 0) {
        $q.notify({ type: 'warning', message: `${unmarked.length} alunni senza presenza assegnata.` })
        return
    }
    saving.value = true
    try {
        const classId = typeof selectedClass.value === 'object' ? selectedClass.value?.id : selectedClass.value

        // Effective subject ID: only pass a valid UUID, never empty string
        const effectiveSubjectId = (!isSubstitutionMode.value && lessonSubjectId.value) ? lessonSubjectId.value : ''

        // 1. Save Attendance Record
        await api.post('/attendance/mark-bulk', {
            class_id: classId,
            date: date.value,
            hour: selectedHour.value,
            subject_id: effectiveSubjectId,
            is_substitution: isSubstitutionMode.value,
            statuses: students.value.map(s => ({
                student_id: s.id,
                status: s.status,
                entry_time: (s.status === 'Late' && s.entry_time?.trim()) ? s.entry_time : null,
                exit_time: (s.status === 'LeftEarly' && s.exit_time?.trim()) ? s.exit_time : null
            }))
        })

        // 2. Save Lesson Signature if topic specified
        if (lessonTopic.value?.trim()) {
            const lessonPayload = {
                class_id: classId,
                subject_id: effectiveSubjectId || null,
                date: date.value,
                hour: selectedHour.value,
                duration: 1,
                topic: lessonTopic.value,
                type: isSubstitutionMode.value ? 'Supplenza' : lessonType.value,
                activity_type: activityType.value || 'standard',
                is_co_teaching: isCoTeaching.value,
                notes: lessonNotes.value
            }
            
            let lessonId = currentHourLesson.value?.id
            if (currentHourLesson.value && !isReadOnly.value) {
                await lessonService.updateLesson(lessonId, lessonPayload)
            } else {
                const lessonRes = await lessonService.createLesson(lessonPayload)
                lessonId = lessonRes.data?.id
            }

            // 3. Save Homework if checked
            if (assignHomework.value && homeworkDesc.value && lessonId) {
                await lessonService.createHomework({
                    class_id: classId,
                    subject_id: effectiveSubjectId || null,
                    lesson_id: lessonId,
                    due_date: homeworkDue.value || date.value,
                    description: homeworkDesc.value
                })
            }
        }

        localStorage.removeItem(`attendance_draft_${classId}_${date.value}_${selectedHour.value}`)
        lastAutosaveTime.value = ''

        $q.notify({ type: 'positive', message: `✓ Registro e Firma Lezione salvati — Ora ${selectedHour.value}` })
        fetchData()

    } catch (error) {
        $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' })
    } finally {
        saving.value = false
    }
}

const deleteUnifiedRecord = () => {
    $q.dialog({
        title: 'Cancella Firma e Presenze',
        message: `Sei sicuro di voler cancellare la firma della lezione e tutte le presenze per la ${selectedHour.value}ª ora del ${date.value}?`,
        cancel: true,
        persistent: true,
        ok: {
            label: 'Elimina',
            color: 'negative'
        }
    }).onOk(async () => {
        saving.value = true
        try {
            const classId = typeof selectedClass.value === 'object' ? selectedClass.value?.id : selectedClass.value

            // 1. Delete lesson if it exists
            if (currentHourLesson.value) {
                await lessonService.deleteLesson(currentHourLesson.value.id)
            }

            // 2. Delete attendance records for this hour
            await attendanceService.deleteHour(classId, date.value, selectedHour.value)

            // 3. Clean local draft if any
            localStorage.removeItem(`attendance_draft_${classId}_${date.value}_${selectedHour.value}`)
            lastAutosaveTime.value = ''

            // 4. Reset form fields for lesson
            lessonTopic.value = ''
            lessonNotes.value = ''
            assignHomework.value = false
            homeworkDesc.value = ''
            homeworkDue.value = ''

            $q.notify({ type: 'positive', message: `✓ Firma e presenze cancellate per l'Ora ${selectedHour.value}` })
            await fetchData()
        } catch (error) {
            console.error(error)
            $q.notify({ type: 'negative', message: 'Errore durante la cancellazione della firma e presenze' })
        } finally {
            saving.value = false
        }
    })
}

const openNoteDialog = (student) => {
    selectedStudentForNote.value = student
    showNoteDialog.value = true
}

const exportCSV = async () => {
    if (!selectedClass.value) {
        $q.notify({ type: 'warning', message: 'Seleziona una classe prima di esportare' })
        return
    }
    try {
        const classId = typeof selectedClass.value === 'object' ? selectedClass.value.id : selectedClass.value
        const res = await api.get('/attendance/export', { params: { class_id: classId, date: date.value }, responseType: 'blob' })
        const url = window.URL.createObjectURL(new Blob([res.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `presenze_${classId}.csv`)
        document.body.appendChild(link)
        link.click()
        link.remove()
        window.URL.revokeObjectURL(url)
        $q.notify({ type: 'positive', message: 'Export CSV completato!' })
    } catch (err) {
        $q.notify({ type: 'negative', message: "Errore durante l'export CSV" })
    }
}
</script>

<style scoped>
.transition-bg {
    transition: background-color 0.3s ease;
}
.shadow-soft {
    box-shadow: 0 4px 12px rgba(0,0,0,0.05);
}
</style>
