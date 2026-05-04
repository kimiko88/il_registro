/// <reference types="../../../../../../../.npm/_npx/2db181330ea4b15b/node_modules/@vue/language-core/types/template-helpers.d.ts" />
/// <reference types="../../../../../../../.npm/_npx/2db181330ea4b15b/node_modules/@vue/language-core/types/props-fallback.d.ts" />
import { ref, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import api from 'src/services/api';
const $q = useQuasar();
const loading = ref(false);
const saving = ref(false);
const activeTab = ref('active');
const projects = ref([]);
const companies = ref([]);
const dialog = ref(false);
const isEdit = ref(false);
const showCompanies = ref(false);
const showAddCompany = ref(false);
const savingCompany = ref(false);
const pctoStats = ref({
    total_projects: 0,
    total_students: 0,
    total_hours: 0,
    active_companies: 0
});
const stats = computed(() => [
    { label: 'Progetti Totali', value: pctoStats.value.total_projects, icon: 'work', color: 'indigo', trendText: 'Dati in tempo reale' },
    { label: 'Studenti Coinvolti', value: pctoStats.value.total_students, icon: 'people', color: 'blue', trendText: 'Iscritti ai percorsi' },
    { label: 'Ore Registrate', value: pctoStats.value.total_hours.toFixed(0), icon: 'timer', color: 'emerald', trendText: 'Ore totali validate' },
    { label: 'Aziende Partner', value: pctoStats.value.active_companies, icon: 'business', color: 'orange', trendText: 'Convenzioni attive' },
]);
const form = ref({
    id: null,
    title: '',
    description: '',
    type: 'Esterno',
    start_date: '',
    end_date: '',
    total_hours: 40,
    company_id: null,
    company_tutor_name: ''
});
const companyForm = ref({
    name: '',
    vat_number: '',
    address: '',
    contact_person: '',
    email: ''
});
const columns = [
    { name: 'title', label: 'Progetto', field: 'title', align: 'left', sortable: true, classes: 'text-weight-bold text-slate-800' },
    { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
    { name: 'hours', label: 'Ore', field: 'total_hours', align: 'center' },
    { name: 'dates', label: 'Periodo', field: row => `${row.start_date.split('T')[0]} - ${row.end_date.split('T')[0]}`, align: 'left' },
    { name: 'status', label: 'Stato', align: 'center' },
    { name: 'actions', label: '', align: 'right' }
];
const filteredProjects = computed(() => {
    const now = new Date();
    if (activeTab.value === 'active') {
        return projects.value.filter(p => new Date(p.end_date) >= now);
    }
    return projects.value.filter(p => new Date(p.end_date) < now);
});
onMounted(async () => {
    await fetchData();
});
const fetchData = async () => {
    loading.value = true;
    try {
        const [pRes, cRes, sRes] = await Promise.all([
            api.get('/pcto/projects'),
            api.get('/pcto/companies'),
            api.get('/pcto/stats')
        ]);
        projects.value = pRes.data || [];
        companies.value = cRes.data || [];
        pctoStats.value = sRes.data || pctoStats.value;
    }
    catch (err) {
        console.error(err);
        $q.notify({ type: 'negative', message: 'Errore durante il caricamento dei dati' });
    }
    finally {
        loading.value = false;
    }
};
const openCreateDialog = () => {
    isEdit.value = false;
    form.value = {
        id: null, title: '', description: '', type: 'Esterno',
        start_date: '', end_date: '', total_hours: 40,
        company_id: null, company_tutor_name: ''
    };
    dialog.value = true;
};
const editProject = (row) => {
    isEdit.value = true;
    form.value = {
        ...row,
        start_date: row.start_date.split('T')[0],
        end_date: row.end_date.split('T')[0]
    };
    dialog.value = true;
};
const saveProject = async () => {
    saving.value = true;
    try {
        if (isEdit.value) {
            await api.put(`/pcto/projects/${form.value.id}`, form.value);
            $q.notify({ type: 'positive', message: 'Progetto aggiornato' });
        }
        else {
            await api.post('/pcto/projects', form.value);
            $q.notify({ type: 'positive', message: 'Progetto creato' });
        }
        dialog.value = false;
        await fetchData();
    }
    catch (err) {
        $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' });
    }
    finally {
        saving.value = false;
    }
};
const confirmDelete = (row) => {
    $q.dialog({
        title: 'Conferma Eliminazione',
        message: `Sei sicuro di voler eliminare il progetto "${row.title}"? Questa azione è irreversibile.`,
        cancel: true,
        persistent: true,
        ok: { color: 'negative', label: 'Elimina', flat: false }
    }).onOk(async () => {
        try {
            await api.delete(`/pcto/projects/${row.id}`);
            $q.notify({ type: 'positive', message: 'Progetto eliminato' });
            await fetchData();
        }
        catch (err) {
            $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione' });
        }
    });
};
const addCompany = async () => {
    savingCompany.value = true;
    try {
        await api.post('/pcto/companies', companyForm.value);
        $q.notify({ type: 'positive', message: 'Azienda aggiunta' });
        showAddCompany.value = false;
        await fetchData();
    }
    catch (err) {
        $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' });
    }
    finally {
        savingCompany.value = false;
    }
};
const getStatusLabel = (row) => {
    const now = new Date();
    const start = new Date(row.start_date);
    const end = new Date(row.end_date);
    if (now < start)
        return 'In arrivo';
    if (now > end)
        return 'Completato';
    return 'In corso';
};
const getStatusColor = (row) => {
    const status = getStatusLabel(row);
    if (status === 'In corso')
        return 'positive';
    if (status === 'Completato')
        return 'grey-7';
    return 'orange-6';
};
const assignStudents = (row) => {
    $q.notify({ message: 'Funzionalità di assegnazione in fase di sviluppo' });
};
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
let __VLS_0;
/** @ts-ignore @type { | typeof __VLS_components.qPage | typeof __VLS_components.QPage | typeof __VLS_components['q-page'] | typeof __VLS_components.qPage | typeof __VLS_components.QPage | typeof __VLS_components['q-page']} */
qPage;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({
    padding: true,
    ...{ class: "bg-slate-50" },
}));
const __VLS_2 = __VLS_1({
    padding: true,
    ...{ class: "bg-slate-50" },
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
var __VLS_5 = {};
/** @type {__VLS_StyleScopedClasses['bg-slate-50']} */ ;
const { default: __VLS_6 } = __VLS_3.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "row items-center justify-between q-mb-xl" },
});
/** @type {__VLS_StyleScopedClasses['row']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
/** @type {__VLS_StyleScopedClasses['q-mb-xl']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({
    ...{ class: "text-h4 text-weight-bold text-outfit q-my-none text-gradient-premium" },
});
/** @type {__VLS_StyleScopedClasses['text-h4']} */ ;
/** @type {__VLS_StyleScopedClasses['text-weight-bold']} */ ;
/** @type {__VLS_StyleScopedClasses['text-outfit']} */ ;
/** @type {__VLS_StyleScopedClasses['q-my-none']} */ ;
/** @type {__VLS_StyleScopedClasses['text-gradient-premium']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
    ...{ class: "text-subtitle1 text-slate-500 q-mt-sm q-mb-none" },
});
/** @type {__VLS_StyleScopedClasses['text-subtitle1']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-500']} */ ;
/** @type {__VLS_StyleScopedClasses['q-mt-sm']} */ ;
/** @type {__VLS_StyleScopedClasses['q-mb-none']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "row items-center q-gutter-sm" },
});
/** @type {__VLS_StyleScopedClasses['row']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['q-gutter-sm']} */ ;
let __VLS_7;
/** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
qBtn;
// @ts-ignore
const __VLS_8 = __VLS_asFunctionalComponent1(__VLS_7, new __VLS_7({
    ...{ 'onClick': {} },
    color: "primary",
    unelevated: true,
    icon: "add",
    label: "Nuovo Progetto",
    ...{ class: "rounded-lg q-px-md shadow-sm" },
    noCaps: true,
}));
const __VLS_9 = __VLS_8({
    ...{ 'onClick': {} },
    color: "primary",
    unelevated: true,
    icon: "add",
    label: "Nuovo Progetto",
    ...{ class: "rounded-lg q-px-md shadow-sm" },
    noCaps: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_8));
let __VLS_12;
const __VLS_13 = ({ click: {} },
    { onClick: (__VLS_ctx.openCreateDialog) });
/** @type {__VLS_StyleScopedClasses['rounded-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['q-px-md']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-sm']} */ ;
var __VLS_10;
var __VLS_11;
let __VLS_14;
/** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
qBtn;
// @ts-ignore
const __VLS_15 = __VLS_asFunctionalComponent1(__VLS_14, new __VLS_14({
    ...{ 'onClick': {} },
    outline: true,
    color: "primary",
    icon: "business",
    label: "Aziende",
    ...{ class: "rounded-lg q-px-md" },
    noCaps: true,
}));
const __VLS_16 = __VLS_15({
    ...{ 'onClick': {} },
    outline: true,
    color: "primary",
    icon: "business",
    label: "Aziende",
    ...{ class: "rounded-lg q-px-md" },
    noCaps: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_15));
let __VLS_19;
const __VLS_20 = ({ click: {} },
    { onClick: (...[$event]) => {
            __VLS_ctx.showCompanies = true;
            // @ts-ignore
            [openCreateDialog, showCompanies,];
        } });
/** @type {__VLS_StyleScopedClasses['rounded-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['q-px-md']} */ ;
var __VLS_17;
var __VLS_18;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "row q-col-gutter-lg q-mb-xl" },
});
/** @type {__VLS_StyleScopedClasses['row']} */ ;
/** @type {__VLS_StyleScopedClasses['q-col-gutter-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['q-mb-xl']} */ ;
for (const [stat] of __VLS_vFor((__VLS_ctx.stats))) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "col-12 col-sm-6 col-md-3" },
        key: (stat.label),
    });
    /** @type {__VLS_StyleScopedClasses['col-12']} */ ;
    /** @type {__VLS_StyleScopedClasses['col-sm-6']} */ ;
    /** @type {__VLS_StyleScopedClasses['col-md-3']} */ ;
    let __VLS_21;
    /** @ts-ignore @type { | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card'] | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card']} */
    qCard;
    // @ts-ignore
    const __VLS_22 = __VLS_asFunctionalComponent1(__VLS_21, new __VLS_21({
        ...{ class: "rounded-xl shadow-soft border-slate-100 overflow-hidden h-full bg-white" },
    }));
    const __VLS_23 = __VLS_22({
        ...{ class: "rounded-xl shadow-soft border-slate-100 overflow-hidden h-full bg-white" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_22));
    /** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
    /** @type {__VLS_StyleScopedClasses['shadow-soft']} */ ;
    /** @type {__VLS_StyleScopedClasses['border-slate-100']} */ ;
    /** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
    /** @type {__VLS_StyleScopedClasses['h-full']} */ ;
    /** @type {__VLS_StyleScopedClasses['bg-white']} */ ;
    const { default: __VLS_26 } = __VLS_24.slots;
    let __VLS_27;
    /** @ts-ignore @type { | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section'] | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section']} */
    qCardSection;
    // @ts-ignore
    const __VLS_28 = __VLS_asFunctionalComponent1(__VLS_27, new __VLS_27({
        ...{ class: "q-pa-lg" },
    }));
    const __VLS_29 = __VLS_28({
        ...{ class: "q-pa-lg" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_28));
    /** @type {__VLS_StyleScopedClasses['q-pa-lg']} */ ;
    const { default: __VLS_32 } = __VLS_30.slots;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "row items-center justify-between q-mb-md" },
    });
    /** @type {__VLS_StyleScopedClasses['row']} */ ;
    /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
    /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-mb-md']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "text-overline text-slate-400 letter-spacing-1" },
    });
    /** @type {__VLS_StyleScopedClasses['text-overline']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
    /** @type {__VLS_StyleScopedClasses['letter-spacing-1']} */ ;
    (stat.label);
    let __VLS_33;
    /** @ts-ignore @type { | typeof __VLS_components.qAvatar | typeof __VLS_components.QAvatar | typeof __VLS_components['q-avatar']} */
    qAvatar;
    // @ts-ignore
    const __VLS_34 = __VLS_asFunctionalComponent1(__VLS_33, new __VLS_33({
        color: (stat.color + '-50'),
        textColor: (stat.color + '-700'),
        icon: (stat.icon),
        size: "40px",
    }));
    const __VLS_35 = __VLS_34({
        color: (stat.color + '-50'),
        textColor: (stat.color + '-700'),
        icon: (stat.icon),
        size: "40px",
    }, ...__VLS_functionalComponentArgsRest(__VLS_34));
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "text-h3 text-weight-bold text-slate-800 q-my-none" },
    });
    /** @type {__VLS_StyleScopedClasses['text-h3']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-weight-bold']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-slate-800']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-my-none']} */ ;
    (stat.value);
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "text-caption text-slate-500 q-mt-sm" },
    });
    /** @type {__VLS_StyleScopedClasses['text-caption']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-slate-500']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-mt-sm']} */ ;
    (stat.trendText);
    // @ts-ignore
    [stats,];
    var __VLS_30;
    // @ts-ignore
    [];
    var __VLS_24;
    // @ts-ignore
    [];
}
let __VLS_38;
/** @ts-ignore @type { | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card'] | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card']} */
qCard;
// @ts-ignore
const __VLS_39 = __VLS_asFunctionalComponent1(__VLS_38, new __VLS_38({
    ...{ class: "rounded-xl shadow-soft border-slate-100 overflow-hidden bg-white" },
}));
const __VLS_40 = __VLS_39({
    ...{ class: "rounded-xl shadow-soft border-slate-100 overflow-hidden bg-white" },
}, ...__VLS_functionalComponentArgsRest(__VLS_39));
/** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-soft']} */ ;
/** @type {__VLS_StyleScopedClasses['border-slate-100']} */ ;
/** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-white']} */ ;
const { default: __VLS_43 } = __VLS_41.slots;
let __VLS_44;
/** @ts-ignore @type { | typeof __VLS_components.qTabs | typeof __VLS_components.QTabs | typeof __VLS_components['q-tabs'] | typeof __VLS_components.qTabs | typeof __VLS_components.QTabs | typeof __VLS_components['q-tabs']} */
qTabs;
// @ts-ignore
const __VLS_45 = __VLS_asFunctionalComponent1(__VLS_44, new __VLS_44({
    modelValue: (__VLS_ctx.activeTab),
    dense: true,
    ...{ class: "text-slate-500 border-b border-slate-100" },
    activeColor: "primary",
    indicatorColor: "primary",
    align: "left",
    narrowIndicator: true,
    noCaps: true,
}));
const __VLS_46 = __VLS_45({
    modelValue: (__VLS_ctx.activeTab),
    dense: true,
    ...{ class: "text-slate-500 border-b border-slate-100" },
    activeColor: "primary",
    indicatorColor: "primary",
    align: "left",
    narrowIndicator: true,
    noCaps: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_45));
/** @type {__VLS_StyleScopedClasses['text-slate-500']} */ ;
/** @type {__VLS_StyleScopedClasses['border-b']} */ ;
/** @type {__VLS_StyleScopedClasses['border-slate-100']} */ ;
const { default: __VLS_49 } = __VLS_47.slots;
let __VLS_50;
/** @ts-ignore @type { | typeof __VLS_components.qTab | typeof __VLS_components.QTab | typeof __VLS_components['q-tab']} */
qTab;
// @ts-ignore
const __VLS_51 = __VLS_asFunctionalComponent1(__VLS_50, new __VLS_50({
    name: "active",
    label: "Progetti Attivi",
    ...{ class: "q-px-lg py-4" },
}));
const __VLS_52 = __VLS_51({
    name: "active",
    label: "Progetti Attivi",
    ...{ class: "q-px-lg py-4" },
}, ...__VLS_functionalComponentArgsRest(__VLS_51));
/** @type {__VLS_StyleScopedClasses['q-px-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['py-4']} */ ;
let __VLS_55;
/** @ts-ignore @type { | typeof __VLS_components.qTab | typeof __VLS_components.QTab | typeof __VLS_components['q-tab']} */
qTab;
// @ts-ignore
const __VLS_56 = __VLS_asFunctionalComponent1(__VLS_55, new __VLS_55({
    name: "archived",
    label: "Archivio Storico",
    ...{ class: "q-px-lg py-4" },
}));
const __VLS_57 = __VLS_56({
    name: "archived",
    label: "Archivio Storico",
    ...{ class: "q-px-lg py-4" },
}, ...__VLS_functionalComponentArgsRest(__VLS_56));
/** @type {__VLS_StyleScopedClasses['q-px-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['py-4']} */ ;
// @ts-ignore
[activeTab,];
var __VLS_47;
let __VLS_60;
/** @ts-ignore @type { | typeof __VLS_components.qTabPanels | typeof __VLS_components.QTabPanels | typeof __VLS_components['q-tab-panels'] | typeof __VLS_components.qTabPanels | typeof __VLS_components.QTabPanels | typeof __VLS_components['q-tab-panels']} */
qTabPanels;
// @ts-ignore
const __VLS_61 = __VLS_asFunctionalComponent1(__VLS_60, new __VLS_60({
    modelValue: (__VLS_ctx.activeTab),
    animated: true,
    ...{ class: "bg-transparent" },
}));
const __VLS_62 = __VLS_61({
    modelValue: (__VLS_ctx.activeTab),
    animated: true,
    ...{ class: "bg-transparent" },
}, ...__VLS_functionalComponentArgsRest(__VLS_61));
/** @type {__VLS_StyleScopedClasses['bg-transparent']} */ ;
const { default: __VLS_65 } = __VLS_63.slots;
let __VLS_66;
/** @ts-ignore @type { | typeof __VLS_components.qTabPanel | typeof __VLS_components.QTabPanel | typeof __VLS_components['q-tab-panel'] | typeof __VLS_components.qTabPanel | typeof __VLS_components.QTabPanel | typeof __VLS_components['q-tab-panel']} */
qTabPanel;
// @ts-ignore
const __VLS_67 = __VLS_asFunctionalComponent1(__VLS_66, new __VLS_66({
    name: "active",
    ...{ class: "q-pa-none" },
}));
const __VLS_68 = __VLS_67({
    name: "active",
    ...{ class: "q-pa-none" },
}, ...__VLS_functionalComponentArgsRest(__VLS_67));
/** @type {__VLS_StyleScopedClasses['q-pa-none']} */ ;
const { default: __VLS_71 } = __VLS_69.slots;
let __VLS_72;
/** @ts-ignore @type { | typeof __VLS_components.qTable | typeof __VLS_components.QTable | typeof __VLS_components['q-table'] | typeof __VLS_components.qTable | typeof __VLS_components.QTable | typeof __VLS_components['q-table']} */
qTable;
// @ts-ignore
const __VLS_73 = __VLS_asFunctionalComponent1(__VLS_72, new __VLS_72({
    rows: (__VLS_ctx.filteredProjects),
    columns: (__VLS_ctx.columns),
    rowKey: "id",
    flat: true,
    loading: (__VLS_ctx.loading),
    ...{ class: "bg-transparent" },
    pagination: ({ rowsPerPage: 10 }),
}));
const __VLS_74 = __VLS_73({
    rows: (__VLS_ctx.filteredProjects),
    columns: (__VLS_ctx.columns),
    rowKey: "id",
    flat: true,
    loading: (__VLS_ctx.loading),
    ...{ class: "bg-transparent" },
    pagination: ({ rowsPerPage: 10 }),
}, ...__VLS_functionalComponentArgsRest(__VLS_73));
/** @type {__VLS_StyleScopedClasses['bg-transparent']} */ ;
const { default: __VLS_77 } = __VLS_75.slots;
{
    const { 'body-cell-type': __VLS_78 } = __VLS_75.slots;
    const [props] = __VLS_vSlot(__VLS_78);
    let __VLS_79;
    /** @ts-ignore @type { | typeof __VLS_components.qTd | typeof __VLS_components.QTd | typeof __VLS_components['q-td'] | typeof __VLS_components.qTd | typeof __VLS_components.QTd | typeof __VLS_components['q-td']} */
    qTd;
    // @ts-ignore
    const __VLS_80 = __VLS_asFunctionalComponent1(__VLS_79, new __VLS_79({
        props: (props),
    }));
    const __VLS_81 = __VLS_80({
        props: (props),
    }, ...__VLS_functionalComponentArgsRest(__VLS_80));
    const { default: __VLS_84 } = __VLS_82.slots;
    let __VLS_85;
    /** @ts-ignore @type { | typeof __VLS_components.qChip | typeof __VLS_components.QChip | typeof __VLS_components['q-chip'] | typeof __VLS_components.qChip | typeof __VLS_components.QChip | typeof __VLS_components['q-chip']} */
    qChip;
    // @ts-ignore
    const __VLS_86 = __VLS_asFunctionalComponent1(__VLS_85, new __VLS_85({
        color: (props.value === 'Interno' ? 'indigo-50' : 'blue-50'),
        textColor: (props.value === 'Interno' ? 'indigo-700' : 'blue-700'),
        size: "sm",
        ...{ class: "text-weight-bold rounded-md" },
    }));
    const __VLS_87 = __VLS_86({
        color: (props.value === 'Interno' ? 'indigo-50' : 'blue-50'),
        textColor: (props.value === 'Interno' ? 'indigo-700' : 'blue-700'),
        size: "sm",
        ...{ class: "text-weight-bold rounded-md" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_86));
    /** @type {__VLS_StyleScopedClasses['text-weight-bold']} */ ;
    /** @type {__VLS_StyleScopedClasses['rounded-md']} */ ;
    const { default: __VLS_90 } = __VLS_88.slots;
    (props.value);
    // @ts-ignore
    [activeTab, filteredProjects, columns, loading,];
    var __VLS_88;
    // @ts-ignore
    [];
    var __VLS_82;
    // @ts-ignore
    [];
}
{
    const { 'body-cell-status': __VLS_91 } = __VLS_75.slots;
    const [props] = __VLS_vSlot(__VLS_91);
    let __VLS_92;
    /** @ts-ignore @type { | typeof __VLS_components.qTd | typeof __VLS_components.QTd | typeof __VLS_components['q-td'] | typeof __VLS_components.qTd | typeof __VLS_components.QTd | typeof __VLS_components['q-td']} */
    qTd;
    // @ts-ignore
    const __VLS_93 = __VLS_asFunctionalComponent1(__VLS_92, new __VLS_92({
        props: (props),
    }));
    const __VLS_94 = __VLS_93({
        props: (props),
    }, ...__VLS_functionalComponentArgsRest(__VLS_93));
    const { default: __VLS_97 } = __VLS_95.slots;
    let __VLS_98;
    /** @ts-ignore @type { | typeof __VLS_components.qBadge | typeof __VLS_components.QBadge | typeof __VLS_components['q-badge'] | typeof __VLS_components.qBadge | typeof __VLS_components.QBadge | typeof __VLS_components['q-badge']} */
    qBadge;
    // @ts-ignore
    const __VLS_99 = __VLS_asFunctionalComponent1(__VLS_98, new __VLS_98({
        color: (__VLS_ctx.getStatusColor(props.row)),
        rounded: true,
        ...{ class: "q-px-sm q-py-xs shadow-xs" },
    }));
    const __VLS_100 = __VLS_99({
        color: (__VLS_ctx.getStatusColor(props.row)),
        rounded: true,
        ...{ class: "q-px-sm q-py-xs shadow-xs" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_99));
    /** @type {__VLS_StyleScopedClasses['q-px-sm']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-py-xs']} */ ;
    /** @type {__VLS_StyleScopedClasses['shadow-xs']} */ ;
    const { default: __VLS_103 } = __VLS_101.slots;
    (__VLS_ctx.getStatusLabel(props.row));
    // @ts-ignore
    [getStatusColor, getStatusLabel,];
    var __VLS_101;
    // @ts-ignore
    [];
    var __VLS_95;
    // @ts-ignore
    [];
}
{
    const { 'body-cell-actions': __VLS_104 } = __VLS_75.slots;
    const [props] = __VLS_vSlot(__VLS_104);
    let __VLS_105;
    /** @ts-ignore @type { | typeof __VLS_components.qTd | typeof __VLS_components.QTd | typeof __VLS_components['q-td'] | typeof __VLS_components.qTd | typeof __VLS_components.QTd | typeof __VLS_components['q-td']} */
    qTd;
    // @ts-ignore
    const __VLS_106 = __VLS_asFunctionalComponent1(__VLS_105, new __VLS_105({
        props: (props),
        ...{ class: "text-right" },
    }));
    const __VLS_107 = __VLS_106({
        props: (props),
        ...{ class: "text-right" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_106));
    /** @type {__VLS_StyleScopedClasses['text-right']} */ ;
    const { default: __VLS_110 } = __VLS_108.slots;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "row q-gutter-xs justify-end" },
    });
    /** @type {__VLS_StyleScopedClasses['row']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-gutter-xs']} */ ;
    /** @type {__VLS_StyleScopedClasses['justify-end']} */ ;
    let __VLS_111;
    /** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn'] | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
    qBtn;
    // @ts-ignore
    const __VLS_112 = __VLS_asFunctionalComponent1(__VLS_111, new __VLS_111({
        ...{ 'onClick': {} },
        flat: true,
        round: true,
        dense: true,
        color: "primary",
        icon: "edit",
    }));
    const __VLS_113 = __VLS_112({
        ...{ 'onClick': {} },
        flat: true,
        round: true,
        dense: true,
        color: "primary",
        icon: "edit",
    }, ...__VLS_functionalComponentArgsRest(__VLS_112));
    let __VLS_116;
    const __VLS_117 = ({ click: {} },
        { onClick: (...[$event]) => {
                __VLS_ctx.editProject(props.row);
                // @ts-ignore
                [editProject,];
            } });
    const { default: __VLS_118 } = __VLS_114.slots;
    let __VLS_119;
    /** @ts-ignore @type { | typeof __VLS_components.qTooltip | typeof __VLS_components.QTooltip | typeof __VLS_components['q-tooltip'] | typeof __VLS_components.qTooltip | typeof __VLS_components.QTooltip | typeof __VLS_components['q-tooltip']} */
    qTooltip;
    // @ts-ignore
    const __VLS_120 = __VLS_asFunctionalComponent1(__VLS_119, new __VLS_119({}));
    const __VLS_121 = __VLS_120({}, ...__VLS_functionalComponentArgsRest(__VLS_120));
    const { default: __VLS_124 } = __VLS_122.slots;
    // @ts-ignore
    [];
    var __VLS_122;
    // @ts-ignore
    [];
    var __VLS_114;
    var __VLS_115;
    let __VLS_125;
    /** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn'] | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
    qBtn;
    // @ts-ignore
    const __VLS_126 = __VLS_asFunctionalComponent1(__VLS_125, new __VLS_125({
        ...{ 'onClick': {} },
        flat: true,
        round: true,
        dense: true,
        color: "indigo",
        icon: "group_add",
    }));
    const __VLS_127 = __VLS_126({
        ...{ 'onClick': {} },
        flat: true,
        round: true,
        dense: true,
        color: "indigo",
        icon: "group_add",
    }, ...__VLS_functionalComponentArgsRest(__VLS_126));
    let __VLS_130;
    const __VLS_131 = ({ click: {} },
        { onClick: (...[$event]) => {
                __VLS_ctx.assignStudents(props.row);
                // @ts-ignore
                [assignStudents,];
            } });
    const { default: __VLS_132 } = __VLS_128.slots;
    let __VLS_133;
    /** @ts-ignore @type { | typeof __VLS_components.qTooltip | typeof __VLS_components.QTooltip | typeof __VLS_components['q-tooltip'] | typeof __VLS_components.qTooltip | typeof __VLS_components.QTooltip | typeof __VLS_components['q-tooltip']} */
    qTooltip;
    // @ts-ignore
    const __VLS_134 = __VLS_asFunctionalComponent1(__VLS_133, new __VLS_133({}));
    const __VLS_135 = __VLS_134({}, ...__VLS_functionalComponentArgsRest(__VLS_134));
    const { default: __VLS_138 } = __VLS_136.slots;
    // @ts-ignore
    [];
    var __VLS_136;
    // @ts-ignore
    [];
    var __VLS_128;
    var __VLS_129;
    let __VLS_139;
    /** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn'] | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
    qBtn;
    // @ts-ignore
    const __VLS_140 = __VLS_asFunctionalComponent1(__VLS_139, new __VLS_139({
        ...{ 'onClick': {} },
        flat: true,
        round: true,
        dense: true,
        color: "negative",
        icon: "delete",
    }));
    const __VLS_141 = __VLS_140({
        ...{ 'onClick': {} },
        flat: true,
        round: true,
        dense: true,
        color: "negative",
        icon: "delete",
    }, ...__VLS_functionalComponentArgsRest(__VLS_140));
    let __VLS_144;
    const __VLS_145 = ({ click: {} },
        { onClick: (...[$event]) => {
                __VLS_ctx.confirmDelete(props.row);
                // @ts-ignore
                [confirmDelete,];
            } });
    const { default: __VLS_146 } = __VLS_142.slots;
    let __VLS_147;
    /** @ts-ignore @type { | typeof __VLS_components.qTooltip | typeof __VLS_components.QTooltip | typeof __VLS_components['q-tooltip'] | typeof __VLS_components.qTooltip | typeof __VLS_components.QTooltip | typeof __VLS_components['q-tooltip']} */
    qTooltip;
    // @ts-ignore
    const __VLS_148 = __VLS_asFunctionalComponent1(__VLS_147, new __VLS_147({}));
    const __VLS_149 = __VLS_148({}, ...__VLS_functionalComponentArgsRest(__VLS_148));
    const { default: __VLS_152 } = __VLS_150.slots;
    // @ts-ignore
    [];
    var __VLS_150;
    // @ts-ignore
    [];
    var __VLS_142;
    var __VLS_143;
    // @ts-ignore
    [];
    var __VLS_108;
    // @ts-ignore
    [];
}
{
    const { 'no-data': __VLS_153 } = __VLS_75.slots;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "full-width q-pa-xl text-center text-slate-400" },
    });
    /** @type {__VLS_StyleScopedClasses['full-width']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-pa-xl']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-center']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
    let __VLS_154;
    /** @ts-ignore @type { | typeof __VLS_components.qIcon | typeof __VLS_components.QIcon | typeof __VLS_components['q-icon']} */
    qIcon;
    // @ts-ignore
    const __VLS_155 = __VLS_asFunctionalComponent1(__VLS_154, new __VLS_154({
        name: "work_outline",
        size: "64px",
        ...{ class: "opacity-20 q-mb-md" },
    }));
    const __VLS_156 = __VLS_155({
        name: "work_outline",
        size: "64px",
        ...{ class: "opacity-20 q-mb-md" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_155));
    /** @type {__VLS_StyleScopedClasses['opacity-20']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-mb-md']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "text-h6" },
    });
    /** @type {__VLS_StyleScopedClasses['text-h6']} */ ;
    // @ts-ignore
    [];
}
// @ts-ignore
[];
var __VLS_75;
// @ts-ignore
[];
var __VLS_69;
let __VLS_159;
/** @ts-ignore @type { | typeof __VLS_components.qTabPanel | typeof __VLS_components.QTabPanel | typeof __VLS_components['q-tab-panel'] | typeof __VLS_components.qTabPanel | typeof __VLS_components.QTabPanel | typeof __VLS_components['q-tab-panel']} */
qTabPanel;
// @ts-ignore
const __VLS_160 = __VLS_asFunctionalComponent1(__VLS_159, new __VLS_159({
    name: "archived",
    ...{ class: "q-pa-none" },
}));
const __VLS_161 = __VLS_160({
    name: "archived",
    ...{ class: "q-pa-none" },
}, ...__VLS_functionalComponentArgsRest(__VLS_160));
/** @type {__VLS_StyleScopedClasses['q-pa-none']} */ ;
const { default: __VLS_164 } = __VLS_162.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "q-pa-xl text-center text-slate-400" },
});
/** @type {__VLS_StyleScopedClasses['q-pa-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['text-center']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
let __VLS_165;
/** @ts-ignore @type { | typeof __VLS_components.qIcon | typeof __VLS_components.QIcon | typeof __VLS_components['q-icon']} */
qIcon;
// @ts-ignore
const __VLS_166 = __VLS_asFunctionalComponent1(__VLS_165, new __VLS_165({
    name: "history",
    size: "64px",
    ...{ class: "opacity-20 q-mb-md" },
}));
const __VLS_167 = __VLS_166({
    name: "history",
    size: "64px",
    ...{ class: "opacity-20 q-mb-md" },
}, ...__VLS_functionalComponentArgsRest(__VLS_166));
/** @type {__VLS_StyleScopedClasses['opacity-20']} */ ;
/** @type {__VLS_StyleScopedClasses['q-mb-md']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "text-h6" },
});
/** @type {__VLS_StyleScopedClasses['text-h6']} */ ;
// @ts-ignore
[];
var __VLS_162;
// @ts-ignore
[];
var __VLS_63;
// @ts-ignore
[];
var __VLS_41;
let __VLS_170;
/** @ts-ignore @type { | typeof __VLS_components.qDialog | typeof __VLS_components.QDialog | typeof __VLS_components['q-dialog'] | typeof __VLS_components.qDialog | typeof __VLS_components.QDialog | typeof __VLS_components['q-dialog']} */
qDialog;
// @ts-ignore
const __VLS_171 = __VLS_asFunctionalComponent1(__VLS_170, new __VLS_170({
    modelValue: (__VLS_ctx.dialog),
    persistent: true,
    ...{ class: "premium-dialog" },
}));
const __VLS_172 = __VLS_171({
    modelValue: (__VLS_ctx.dialog),
    persistent: true,
    ...{ class: "premium-dialog" },
}, ...__VLS_functionalComponentArgsRest(__VLS_171));
/** @type {__VLS_StyleScopedClasses['premium-dialog']} */ ;
const { default: __VLS_175 } = __VLS_173.slots;
let __VLS_176;
/** @ts-ignore @type { | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card'] | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card']} */
qCard;
// @ts-ignore
const __VLS_177 = __VLS_asFunctionalComponent1(__VLS_176, new __VLS_176({
    ...{ style: {} },
    ...{ class: "rounded-xl overflow-hidden shadow-24" },
}));
const __VLS_178 = __VLS_177({
    ...{ style: {} },
    ...{ class: "rounded-xl overflow-hidden shadow-24" },
}, ...__VLS_functionalComponentArgsRest(__VLS_177));
/** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-24']} */ ;
const { default: __VLS_181 } = __VLS_179.slots;
let __VLS_182;
/** @ts-ignore @type { | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section'] | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section']} */
qCardSection;
// @ts-ignore
const __VLS_183 = __VLS_asFunctionalComponent1(__VLS_182, new __VLS_182({
    ...{ class: "bg-gradient-primary text-white q-pa-lg" },
}));
const __VLS_184 = __VLS_183({
    ...{ class: "bg-gradient-primary text-white q-pa-lg" },
}, ...__VLS_functionalComponentArgsRest(__VLS_183));
/** @type {__VLS_StyleScopedClasses['bg-gradient-primary']} */ ;
/** @type {__VLS_StyleScopedClasses['text-white']} */ ;
/** @type {__VLS_StyleScopedClasses['q-pa-lg']} */ ;
const { default: __VLS_187 } = __VLS_185.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "text-h5 text-weight-bold" },
});
/** @type {__VLS_StyleScopedClasses['text-h5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-weight-bold']} */ ;
(__VLS_ctx.isEdit ? 'Modifica Progetto' : 'Nuovo Progetto PCTO');
// @ts-ignore
[dialog, isEdit,];
var __VLS_185;
let __VLS_188;
/** @ts-ignore @type { | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section'] | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section']} */
qCardSection;
// @ts-ignore
const __VLS_189 = __VLS_asFunctionalComponent1(__VLS_188, new __VLS_188({
    ...{ class: "q-pa-xl" },
}));
const __VLS_190 = __VLS_189({
    ...{ class: "q-pa-xl" },
}, ...__VLS_functionalComponentArgsRest(__VLS_189));
/** @type {__VLS_StyleScopedClasses['q-pa-xl']} */ ;
const { default: __VLS_193 } = __VLS_191.slots;
let __VLS_194;
/** @ts-ignore @type { | typeof __VLS_components.qForm | typeof __VLS_components.QForm | typeof __VLS_components['q-form'] | typeof __VLS_components.qForm | typeof __VLS_components.QForm | typeof __VLS_components['q-form']} */
qForm;
// @ts-ignore
const __VLS_195 = __VLS_asFunctionalComponent1(__VLS_194, new __VLS_194({
    ...{ 'onSubmit': {} },
    ...{ class: "q-gutter-y-lg" },
}));
const __VLS_196 = __VLS_195({
    ...{ 'onSubmit': {} },
    ...{ class: "q-gutter-y-lg" },
}, ...__VLS_functionalComponentArgsRest(__VLS_195));
let __VLS_199;
const __VLS_200 = ({ submit: {} },
    { onSubmit: (__VLS_ctx.saveProject) });
/** @type {__VLS_StyleScopedClasses['q-gutter-y-lg']} */ ;
const { default: __VLS_201 } = __VLS_197.slots;
let __VLS_202;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_203 = __VLS_asFunctionalComponent1(__VLS_202, new __VLS_202({
    modelValue: (__VLS_ctx.form.title),
    label: "Titolo Progetto",
    outlined: true,
    rules: ([val => !!val || 'Campo richiesto']),
}));
const __VLS_204 = __VLS_203({
    modelValue: (__VLS_ctx.form.title),
    label: "Titolo Progetto",
    outlined: true,
    rules: ([val => !!val || 'Campo richiesto']),
}, ...__VLS_functionalComponentArgsRest(__VLS_203));
let __VLS_207;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_208 = __VLS_asFunctionalComponent1(__VLS_207, new __VLS_207({
    modelValue: (__VLS_ctx.form.description),
    type: "textarea",
    label: "Descrizione",
    outlined: true,
    rows: "3",
}));
const __VLS_209 = __VLS_208({
    modelValue: (__VLS_ctx.form.description),
    type: "textarea",
    label: "Descrizione",
    outlined: true,
    rows: "3",
}, ...__VLS_functionalComponentArgsRest(__VLS_208));
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "row q-col-gutter-lg" },
});
/** @type {__VLS_StyleScopedClasses['row']} */ ;
/** @type {__VLS_StyleScopedClasses['q-col-gutter-lg']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "col-6" },
});
/** @type {__VLS_StyleScopedClasses['col-6']} */ ;
let __VLS_212;
/** @ts-ignore @type { | typeof __VLS_components.qSelect | typeof __VLS_components.QSelect | typeof __VLS_components['q-select']} */
qSelect;
// @ts-ignore
const __VLS_213 = __VLS_asFunctionalComponent1(__VLS_212, new __VLS_212({
    modelValue: (__VLS_ctx.form.type),
    options: (['Interno', 'Esterno']),
    label: "Tipologia",
    outlined: true,
}));
const __VLS_214 = __VLS_213({
    modelValue: (__VLS_ctx.form.type),
    options: (['Interno', 'Esterno']),
    label: "Tipologia",
    outlined: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_213));
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "col-6" },
});
/** @type {__VLS_StyleScopedClasses['col-6']} */ ;
let __VLS_217;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_218 = __VLS_asFunctionalComponent1(__VLS_217, new __VLS_217({
    modelValue: (__VLS_ctx.form.total_hours),
    modelModifiers: { number: true, },
    type: "number",
    label: "Ore Totali",
    outlined: true,
}));
const __VLS_219 = __VLS_218({
    modelValue: (__VLS_ctx.form.total_hours),
    modelModifiers: { number: true, },
    type: "number",
    label: "Ore Totali",
    outlined: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_218));
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "row q-col-gutter-lg" },
});
/** @type {__VLS_StyleScopedClasses['row']} */ ;
/** @type {__VLS_StyleScopedClasses['q-col-gutter-lg']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "col-6" },
});
/** @type {__VLS_StyleScopedClasses['col-6']} */ ;
let __VLS_222;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_223 = __VLS_asFunctionalComponent1(__VLS_222, new __VLS_222({
    modelValue: (__VLS_ctx.form.start_date),
    type: "date",
    label: "Data Inizio",
    outlined: true,
    stackLabel: true,
}));
const __VLS_224 = __VLS_223({
    modelValue: (__VLS_ctx.form.start_date),
    type: "date",
    label: "Data Inizio",
    outlined: true,
    stackLabel: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_223));
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "col-6" },
});
/** @type {__VLS_StyleScopedClasses['col-6']} */ ;
let __VLS_227;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_228 = __VLS_asFunctionalComponent1(__VLS_227, new __VLS_227({
    modelValue: (__VLS_ctx.form.end_date),
    type: "date",
    label: "Data Fine",
    outlined: true,
    stackLabel: true,
}));
const __VLS_229 = __VLS_228({
    modelValue: (__VLS_ctx.form.end_date),
    type: "date",
    label: "Data Fine",
    outlined: true,
    stackLabel: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_228));
if (__VLS_ctx.form.type === 'Esterno') {
    let __VLS_232;
    /** @ts-ignore @type { | typeof __VLS_components.qSelect | typeof __VLS_components.QSelect | typeof __VLS_components['q-select']} */
    qSelect;
    // @ts-ignore
    const __VLS_233 = __VLS_asFunctionalComponent1(__VLS_232, new __VLS_232({
        modelValue: (__VLS_ctx.form.company_id),
        options: (__VLS_ctx.companies),
        optionLabel: "name",
        optionValue: "id",
        emitValue: true,
        mapOptions: true,
        label: "Azienda Ospitante",
        outlined: true,
    }));
    const __VLS_234 = __VLS_233({
        modelValue: (__VLS_ctx.form.company_id),
        options: (__VLS_ctx.companies),
        optionLabel: "name",
        optionValue: "id",
        emitValue: true,
        mapOptions: true,
        label: "Azienda Ospitante",
        outlined: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_233));
}
let __VLS_237;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_238 = __VLS_asFunctionalComponent1(__VLS_237, new __VLS_237({
    modelValue: (__VLS_ctx.form.company_tutor_name),
    label: "Tutor Aziendale",
    outlined: true,
}));
const __VLS_239 = __VLS_238({
    modelValue: (__VLS_ctx.form.company_tutor_name),
    label: "Tutor Aziendale",
    outlined: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_238));
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "row justify-end q-mt-xl q-gutter-sm" },
});
/** @type {__VLS_StyleScopedClasses['row']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-end']} */ ;
/** @type {__VLS_StyleScopedClasses['q-mt-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['q-gutter-sm']} */ ;
let __VLS_242;
/** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
qBtn;
// @ts-ignore
const __VLS_243 = __VLS_asFunctionalComponent1(__VLS_242, new __VLS_242({
    flat: true,
    label: "Annulla",
    color: "slate-400",
    noCaps: true,
}));
const __VLS_244 = __VLS_243({
    flat: true,
    label: "Annulla",
    color: "slate-400",
    noCaps: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_243));
__VLS_asFunctionalDirective(__VLS_directives.vClosePopup, {})(null, { ...__VLS_directiveBindingRestFields, }, null, null);
let __VLS_247;
/** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
qBtn;
// @ts-ignore
const __VLS_248 = __VLS_asFunctionalComponent1(__VLS_247, new __VLS_247({
    type: "submit",
    color: "primary",
    label: (__VLS_ctx.isEdit ? 'Aggiorna Progetto' : 'Crea Progetto'),
    ...{ class: "q-px-xl rounded-lg shadow-sm" },
    noCaps: true,
    loading: (__VLS_ctx.saving),
}));
const __VLS_249 = __VLS_248({
    type: "submit",
    color: "primary",
    label: (__VLS_ctx.isEdit ? 'Aggiorna Progetto' : 'Crea Progetto'),
    ...{ class: "q-px-xl rounded-lg shadow-sm" },
    noCaps: true,
    loading: (__VLS_ctx.saving),
}, ...__VLS_functionalComponentArgsRest(__VLS_248));
/** @type {__VLS_StyleScopedClasses['q-px-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-sm']} */ ;
// @ts-ignore
[isEdit, saveProject, form, form, form, form, form, form, form, form, form, companies, vClosePopup, saving,];
var __VLS_197;
var __VLS_198;
// @ts-ignore
[];
var __VLS_191;
// @ts-ignore
[];
var __VLS_179;
// @ts-ignore
[];
var __VLS_173;
let __VLS_252;
/** @ts-ignore @type { | typeof __VLS_components.qDialog | typeof __VLS_components.QDialog | typeof __VLS_components['q-dialog'] | typeof __VLS_components.qDialog | typeof __VLS_components.QDialog | typeof __VLS_components['q-dialog']} */
qDialog;
// @ts-ignore
const __VLS_253 = __VLS_asFunctionalComponent1(__VLS_252, new __VLS_252({
    modelValue: (__VLS_ctx.showCompanies),
    maximized: true,
    transitionShow: "slide-up",
    transitionHide: "slide-down",
}));
const __VLS_254 = __VLS_253({
    modelValue: (__VLS_ctx.showCompanies),
    maximized: true,
    transitionShow: "slide-up",
    transitionHide: "slide-down",
}, ...__VLS_functionalComponentArgsRest(__VLS_253));
const { default: __VLS_257 } = __VLS_255.slots;
let __VLS_258;
/** @ts-ignore @type { | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card'] | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card']} */
qCard;
// @ts-ignore
const __VLS_259 = __VLS_asFunctionalComponent1(__VLS_258, new __VLS_258({
    ...{ class: "bg-slate-50" },
}));
const __VLS_260 = __VLS_259({
    ...{ class: "bg-slate-50" },
}, ...__VLS_functionalComponentArgsRest(__VLS_259));
/** @type {__VLS_StyleScopedClasses['bg-slate-50']} */ ;
const { default: __VLS_263 } = __VLS_261.slots;
let __VLS_264;
/** @ts-ignore @type { | typeof __VLS_components.qToolbar | typeof __VLS_components.QToolbar | typeof __VLS_components['q-toolbar'] | typeof __VLS_components.qToolbar | typeof __VLS_components.QToolbar | typeof __VLS_components['q-toolbar']} */
qToolbar;
// @ts-ignore
const __VLS_265 = __VLS_asFunctionalComponent1(__VLS_264, new __VLS_264({
    ...{ class: "bg-white border-b border-slate-100 q-py-md" },
}));
const __VLS_266 = __VLS_265({
    ...{ class: "bg-white border-b border-slate-100 q-py-md" },
}, ...__VLS_functionalComponentArgsRest(__VLS_265));
/** @type {__VLS_StyleScopedClasses['bg-white']} */ ;
/** @type {__VLS_StyleScopedClasses['border-b']} */ ;
/** @type {__VLS_StyleScopedClasses['border-slate-100']} */ ;
/** @type {__VLS_StyleScopedClasses['q-py-md']} */ ;
const { default: __VLS_269 } = __VLS_267.slots;
let __VLS_270;
/** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
qBtn;
// @ts-ignore
const __VLS_271 = __VLS_asFunctionalComponent1(__VLS_270, new __VLS_270({
    flat: true,
    round: true,
    dense: true,
    icon: "arrow_back",
    ...{ class: "text-slate-400" },
}));
const __VLS_272 = __VLS_271({
    flat: true,
    round: true,
    dense: true,
    icon: "arrow_back",
    ...{ class: "text-slate-400" },
}, ...__VLS_functionalComponentArgsRest(__VLS_271));
__VLS_asFunctionalDirective(__VLS_directives.vClosePopup, {})(null, { ...__VLS_directiveBindingRestFields, }, null, null);
/** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
let __VLS_275;
/** @ts-ignore @type { | typeof __VLS_components.qToolbarTitle | typeof __VLS_components.QToolbarTitle | typeof __VLS_components['q-toolbar-title'] | typeof __VLS_components.qToolbarTitle | typeof __VLS_components.QToolbarTitle | typeof __VLS_components['q-toolbar-title']} */
qToolbarTitle;
// @ts-ignore
const __VLS_276 = __VLS_asFunctionalComponent1(__VLS_275, new __VLS_275({
    ...{ class: "text-weight-bold text-slate-800" },
}));
const __VLS_277 = __VLS_276({
    ...{ class: "text-weight-bold text-slate-800" },
}, ...__VLS_functionalComponentArgsRest(__VLS_276));
/** @type {__VLS_StyleScopedClasses['text-weight-bold']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-800']} */ ;
const { default: __VLS_280 } = __VLS_278.slots;
// @ts-ignore
[showCompanies, vClosePopup,];
var __VLS_278;
let __VLS_281;
/** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
qBtn;
// @ts-ignore
const __VLS_282 = __VLS_asFunctionalComponent1(__VLS_281, new __VLS_281({
    ...{ 'onClick': {} },
    unelevated: true,
    color: "primary",
    label: "Aggiungi Azienda",
    ...{ class: "rounded-lg q-px-md" },
    noCaps: true,
}));
const __VLS_283 = __VLS_282({
    ...{ 'onClick': {} },
    unelevated: true,
    color: "primary",
    label: "Aggiungi Azienda",
    ...{ class: "rounded-lg q-px-md" },
    noCaps: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_282));
let __VLS_286;
const __VLS_287 = ({ click: {} },
    { onClick: (...[$event]) => {
            __VLS_ctx.showAddCompany = true;
            // @ts-ignore
            [showAddCompany,];
        } });
/** @type {__VLS_StyleScopedClasses['rounded-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['q-px-md']} */ ;
var __VLS_284;
var __VLS_285;
// @ts-ignore
[];
var __VLS_267;
let __VLS_288;
/** @ts-ignore @type { | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section'] | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section']} */
qCardSection;
// @ts-ignore
const __VLS_289 = __VLS_asFunctionalComponent1(__VLS_288, new __VLS_288({
    ...{ class: "q-pa-xl" },
}));
const __VLS_290 = __VLS_289({
    ...{ class: "q-pa-xl" },
}, ...__VLS_functionalComponentArgsRest(__VLS_289));
/** @type {__VLS_StyleScopedClasses['q-pa-xl']} */ ;
const { default: __VLS_293 } = __VLS_291.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "row q-col-gutter-xl" },
});
/** @type {__VLS_StyleScopedClasses['row']} */ ;
/** @type {__VLS_StyleScopedClasses['q-col-gutter-xl']} */ ;
for (const [company] of __VLS_vFor((__VLS_ctx.companies))) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "col-12 col-sm-6 col-md-4" },
        key: (company.id),
    });
    /** @type {__VLS_StyleScopedClasses['col-12']} */ ;
    /** @type {__VLS_StyleScopedClasses['col-sm-6']} */ ;
    /** @type {__VLS_StyleScopedClasses['col-md-4']} */ ;
    let __VLS_294;
    /** @ts-ignore @type { | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card'] | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card']} */
    qCard;
    // @ts-ignore
    const __VLS_295 = __VLS_asFunctionalComponent1(__VLS_294, new __VLS_294({
        flat: true,
        ...{ class: "rounded-xl border-slate-100 bg-white shadow-soft h-full overflow-hidden" },
    }));
    const __VLS_296 = __VLS_295({
        flat: true,
        ...{ class: "rounded-xl border-slate-100 bg-white shadow-soft h-full overflow-hidden" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_295));
    /** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
    /** @type {__VLS_StyleScopedClasses['border-slate-100']} */ ;
    /** @type {__VLS_StyleScopedClasses['bg-white']} */ ;
    /** @type {__VLS_StyleScopedClasses['shadow-soft']} */ ;
    /** @type {__VLS_StyleScopedClasses['h-full']} */ ;
    /** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
    const { default: __VLS_299 } = __VLS_297.slots;
    let __VLS_300;
    /** @ts-ignore @type { | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section'] | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section']} */
    qCardSection;
    // @ts-ignore
    const __VLS_301 = __VLS_asFunctionalComponent1(__VLS_300, new __VLS_300({
        ...{ class: "q-pa-lg" },
    }));
    const __VLS_302 = __VLS_301({
        ...{ class: "q-pa-lg" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_301));
    /** @type {__VLS_StyleScopedClasses['q-pa-lg']} */ ;
    const { default: __VLS_305 } = __VLS_303.slots;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "row items-center no-wrap q-mb-lg" },
    });
    /** @type {__VLS_StyleScopedClasses['row']} */ ;
    /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
    /** @type {__VLS_StyleScopedClasses['no-wrap']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-mb-lg']} */ ;
    let __VLS_306;
    /** @ts-ignore @type { | typeof __VLS_components.qAvatar | typeof __VLS_components.QAvatar | typeof __VLS_components['q-avatar']} */
    qAvatar;
    // @ts-ignore
    const __VLS_307 = __VLS_asFunctionalComponent1(__VLS_306, new __VLS_306({
        color: "blue-50",
        textColor: "blue-700",
        icon: "business",
        size: "48px",
    }));
    const __VLS_308 = __VLS_307({
        color: "blue-50",
        textColor: "blue-700",
        icon: "business",
        size: "48px",
    }, ...__VLS_functionalComponentArgsRest(__VLS_307));
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "q-ml-md" },
    });
    /** @type {__VLS_StyleScopedClasses['q-ml-md']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "text-subtitle1 text-weight-bold text-slate-800" },
    });
    /** @type {__VLS_StyleScopedClasses['text-subtitle1']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-weight-bold']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-slate-800']} */ ;
    (company.name);
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "text-caption text-slate-400" },
    });
    /** @type {__VLS_StyleScopedClasses['text-caption']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
    (company.vat_number);
    let __VLS_311;
    /** @ts-ignore @type { | typeof __VLS_components.qSeparator | typeof __VLS_components.QSeparator | typeof __VLS_components['q-separator']} */
    qSeparator;
    // @ts-ignore
    const __VLS_312 = __VLS_asFunctionalComponent1(__VLS_311, new __VLS_311({
        ...{ class: "q-my-lg opacity-50" },
    }));
    const __VLS_313 = __VLS_312({
        ...{ class: "q-my-lg opacity-50" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_312));
    /** @type {__VLS_StyleScopedClasses['q-my-lg']} */ ;
    /** @type {__VLS_StyleScopedClasses['opacity-50']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "q-gutter-y-sm" },
    });
    /** @type {__VLS_StyleScopedClasses['q-gutter-y-sm']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "row items-center" },
    });
    /** @type {__VLS_StyleScopedClasses['row']} */ ;
    /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
    let __VLS_316;
    /** @ts-ignore @type { | typeof __VLS_components.qIcon | typeof __VLS_components.QIcon | typeof __VLS_components['q-icon']} */
    qIcon;
    // @ts-ignore
    const __VLS_317 = __VLS_asFunctionalComponent1(__VLS_316, new __VLS_316({
        name: "place",
        color: "slate-300",
        size: "18px",
        ...{ class: "q-mr-sm" },
    }));
    const __VLS_318 = __VLS_317({
        name: "place",
        color: "slate-300",
        size: "18px",
        ...{ class: "q-mr-sm" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_317));
    /** @type {__VLS_StyleScopedClasses['q-mr-sm']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
        ...{ class: "text-caption text-slate-600" },
    });
    /** @type {__VLS_StyleScopedClasses['text-caption']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-slate-600']} */ ;
    (company.address);
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "row items-center" },
    });
    /** @type {__VLS_StyleScopedClasses['row']} */ ;
    /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
    let __VLS_321;
    /** @ts-ignore @type { | typeof __VLS_components.qIcon | typeof __VLS_components.QIcon | typeof __VLS_components['q-icon']} */
    qIcon;
    // @ts-ignore
    const __VLS_322 = __VLS_asFunctionalComponent1(__VLS_321, new __VLS_321({
        name: "person",
        color: "slate-300",
        size: "18px",
        ...{ class: "q-mr-sm" },
    }));
    const __VLS_323 = __VLS_322({
        name: "person",
        color: "slate-300",
        size: "18px",
        ...{ class: "q-mr-sm" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_322));
    /** @type {__VLS_StyleScopedClasses['q-mr-sm']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
        ...{ class: "text-caption text-slate-600" },
    });
    /** @type {__VLS_StyleScopedClasses['text-caption']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-slate-600']} */ ;
    (company.contact_person || 'Contatto non definito');
    // @ts-ignore
    [companies,];
    var __VLS_303;
    let __VLS_326;
    /** @ts-ignore @type { | typeof __VLS_components.qCardActions | typeof __VLS_components.QCardActions | typeof __VLS_components['q-card-actions'] | typeof __VLS_components.qCardActions | typeof __VLS_components.QCardActions | typeof __VLS_components['q-card-actions']} */
    qCardActions;
    // @ts-ignore
    const __VLS_327 = __VLS_asFunctionalComponent1(__VLS_326, new __VLS_326({
        align: "right",
        ...{ class: "bg-slate-50 q-pa-md" },
    }));
    const __VLS_328 = __VLS_327({
        align: "right",
        ...{ class: "bg-slate-50 q-pa-md" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_327));
    /** @type {__VLS_StyleScopedClasses['bg-slate-50']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-pa-md']} */ ;
    const { default: __VLS_331 } = __VLS_329.slots;
    let __VLS_332;
    /** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
    qBtn;
    // @ts-ignore
    const __VLS_333 = __VLS_asFunctionalComponent1(__VLS_332, new __VLS_332({
        flat: true,
        dense: true,
        color: "primary",
        icon: "edit",
        label: "Modifica",
        size: "sm",
        noCaps: true,
    }));
    const __VLS_334 = __VLS_333({
        flat: true,
        dense: true,
        color: "primary",
        icon: "edit",
        label: "Modifica",
        size: "sm",
        noCaps: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_333));
    let __VLS_337;
    /** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
    qBtn;
    // @ts-ignore
    const __VLS_338 = __VLS_asFunctionalComponent1(__VLS_337, new __VLS_337({
        flat: true,
        dense: true,
        color: "negative",
        icon: "delete",
        label: "Rimuovi",
        size: "sm",
        noCaps: true,
    }));
    const __VLS_339 = __VLS_338({
        flat: true,
        dense: true,
        color: "negative",
        icon: "delete",
        label: "Rimuovi",
        size: "sm",
        noCaps: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_338));
    // @ts-ignore
    [];
    var __VLS_329;
    // @ts-ignore
    [];
    var __VLS_297;
    // @ts-ignore
    [];
}
if (__VLS_ctx.companies.length === 0) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "full-width q-pa-xl text-center text-slate-400" },
    });
    /** @type {__VLS_StyleScopedClasses['full-width']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-pa-xl']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-center']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
    let __VLS_342;
    /** @ts-ignore @type { | typeof __VLS_components.qIcon | typeof __VLS_components.QIcon | typeof __VLS_components['q-icon']} */
    qIcon;
    // @ts-ignore
    const __VLS_343 = __VLS_asFunctionalComponent1(__VLS_342, new __VLS_342({
        name: "business_center",
        size: "64px",
        ...{ class: "opacity-20 q-mb-md" },
    }));
    const __VLS_344 = __VLS_343({
        name: "business_center",
        size: "64px",
        ...{ class: "opacity-20 q-mb-md" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_343));
    /** @type {__VLS_StyleScopedClasses['opacity-20']} */ ;
    /** @type {__VLS_StyleScopedClasses['q-mb-md']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "text-h6" },
    });
    /** @type {__VLS_StyleScopedClasses['text-h6']} */ ;
}
// @ts-ignore
[companies,];
var __VLS_291;
// @ts-ignore
[];
var __VLS_261;
// @ts-ignore
[];
var __VLS_255;
let __VLS_347;
/** @ts-ignore @type { | typeof __VLS_components.qDialog | typeof __VLS_components.QDialog | typeof __VLS_components['q-dialog'] | typeof __VLS_components.qDialog | typeof __VLS_components.QDialog | typeof __VLS_components['q-dialog']} */
qDialog;
// @ts-ignore
const __VLS_348 = __VLS_asFunctionalComponent1(__VLS_347, new __VLS_347({
    modelValue: (__VLS_ctx.showAddCompany),
    ...{ class: "premium-dialog" },
}));
const __VLS_349 = __VLS_348({
    modelValue: (__VLS_ctx.showAddCompany),
    ...{ class: "premium-dialog" },
}, ...__VLS_functionalComponentArgsRest(__VLS_348));
/** @type {__VLS_StyleScopedClasses['premium-dialog']} */ ;
const { default: __VLS_352 } = __VLS_350.slots;
let __VLS_353;
/** @ts-ignore @type { | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card'] | typeof __VLS_components.qCard | typeof __VLS_components.QCard | typeof __VLS_components['q-card']} */
qCard;
// @ts-ignore
const __VLS_354 = __VLS_asFunctionalComponent1(__VLS_353, new __VLS_353({
    ...{ style: {} },
    ...{ class: "rounded-xl overflow-hidden shadow-24 bg-white" },
}));
const __VLS_355 = __VLS_354({
    ...{ style: {} },
    ...{ class: "rounded-xl overflow-hidden shadow-24 bg-white" },
}, ...__VLS_functionalComponentArgsRest(__VLS_354));
/** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-24']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-white']} */ ;
const { default: __VLS_358 } = __VLS_356.slots;
let __VLS_359;
/** @ts-ignore @type { | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section'] | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section']} */
qCardSection;
// @ts-ignore
const __VLS_360 = __VLS_asFunctionalComponent1(__VLS_359, new __VLS_359({
    ...{ class: "bg-gradient-primary text-white q-pa-lg" },
}));
const __VLS_361 = __VLS_360({
    ...{ class: "bg-gradient-primary text-white q-pa-lg" },
}, ...__VLS_functionalComponentArgsRest(__VLS_360));
/** @type {__VLS_StyleScopedClasses['bg-gradient-primary']} */ ;
/** @type {__VLS_StyleScopedClasses['text-white']} */ ;
/** @type {__VLS_StyleScopedClasses['q-pa-lg']} */ ;
const { default: __VLS_364 } = __VLS_362.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "text-h5 text-weight-bold" },
});
/** @type {__VLS_StyleScopedClasses['text-h5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-weight-bold']} */ ;
// @ts-ignore
[showAddCompany,];
var __VLS_362;
let __VLS_365;
/** @ts-ignore @type { | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section'] | typeof __VLS_components.qCardSection | typeof __VLS_components.QCardSection | typeof __VLS_components['q-card-section']} */
qCardSection;
// @ts-ignore
const __VLS_366 = __VLS_asFunctionalComponent1(__VLS_365, new __VLS_365({
    ...{ class: "q-pa-xl" },
}));
const __VLS_367 = __VLS_366({
    ...{ class: "q-pa-xl" },
}, ...__VLS_functionalComponentArgsRest(__VLS_366));
/** @type {__VLS_StyleScopedClasses['q-pa-xl']} */ ;
const { default: __VLS_370 } = __VLS_368.slots;
let __VLS_371;
/** @ts-ignore @type { | typeof __VLS_components.qForm | typeof __VLS_components.QForm | typeof __VLS_components['q-form'] | typeof __VLS_components.qForm | typeof __VLS_components.QForm | typeof __VLS_components['q-form']} */
qForm;
// @ts-ignore
const __VLS_372 = __VLS_asFunctionalComponent1(__VLS_371, new __VLS_371({
    ...{ 'onSubmit': {} },
    ...{ class: "q-gutter-y-lg" },
}));
const __VLS_373 = __VLS_372({
    ...{ 'onSubmit': {} },
    ...{ class: "q-gutter-y-lg" },
}, ...__VLS_functionalComponentArgsRest(__VLS_372));
let __VLS_376;
const __VLS_377 = ({ submit: {} },
    { onSubmit: (__VLS_ctx.addCompany) });
/** @type {__VLS_StyleScopedClasses['q-gutter-y-lg']} */ ;
const { default: __VLS_378 } = __VLS_374.slots;
let __VLS_379;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_380 = __VLS_asFunctionalComponent1(__VLS_379, new __VLS_379({
    modelValue: (__VLS_ctx.companyForm.name),
    label: "Ragione Sociale",
    outlined: true,
}));
const __VLS_381 = __VLS_380({
    modelValue: (__VLS_ctx.companyForm.name),
    label: "Ragione Sociale",
    outlined: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_380));
let __VLS_384;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_385 = __VLS_asFunctionalComponent1(__VLS_384, new __VLS_384({
    modelValue: (__VLS_ctx.companyForm.vat_number),
    label: "Partita IVA",
    outlined: true,
}));
const __VLS_386 = __VLS_385({
    modelValue: (__VLS_ctx.companyForm.vat_number),
    label: "Partita IVA",
    outlined: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_385));
let __VLS_389;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_390 = __VLS_asFunctionalComponent1(__VLS_389, new __VLS_389({
    modelValue: (__VLS_ctx.companyForm.address),
    label: "Indirizzo Sede",
    outlined: true,
}));
const __VLS_391 = __VLS_390({
    modelValue: (__VLS_ctx.companyForm.address),
    label: "Indirizzo Sede",
    outlined: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_390));
let __VLS_394;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_395 = __VLS_asFunctionalComponent1(__VLS_394, new __VLS_394({
    modelValue: (__VLS_ctx.companyForm.contact_person),
    label: "Persona di Riferimento",
    outlined: true,
}));
const __VLS_396 = __VLS_395({
    modelValue: (__VLS_ctx.companyForm.contact_person),
    label: "Persona di Riferimento",
    outlined: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_395));
let __VLS_399;
/** @ts-ignore @type { | typeof __VLS_components.qInput | typeof __VLS_components.QInput | typeof __VLS_components['q-input']} */
qInput;
// @ts-ignore
const __VLS_400 = __VLS_asFunctionalComponent1(__VLS_399, new __VLS_399({
    modelValue: (__VLS_ctx.companyForm.email),
    label: "Email Contatto",
    outlined: true,
}));
const __VLS_401 = __VLS_400({
    modelValue: (__VLS_ctx.companyForm.email),
    label: "Email Contatto",
    outlined: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_400));
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "row justify-end q-mt-xl q-gutter-sm" },
});
/** @type {__VLS_StyleScopedClasses['row']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-end']} */ ;
/** @type {__VLS_StyleScopedClasses['q-mt-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['q-gutter-sm']} */ ;
let __VLS_404;
/** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
qBtn;
// @ts-ignore
const __VLS_405 = __VLS_asFunctionalComponent1(__VLS_404, new __VLS_404({
    flat: true,
    label: "Annulla",
    color: "slate-400",
    noCaps: true,
}));
const __VLS_406 = __VLS_405({
    flat: true,
    label: "Annulla",
    color: "slate-400",
    noCaps: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_405));
__VLS_asFunctionalDirective(__VLS_directives.vClosePopup, {})(null, { ...__VLS_directiveBindingRestFields, }, null, null);
let __VLS_409;
/** @ts-ignore @type { | typeof __VLS_components.qBtn | typeof __VLS_components.QBtn | typeof __VLS_components['q-btn']} */
qBtn;
// @ts-ignore
const __VLS_410 = __VLS_asFunctionalComponent1(__VLS_409, new __VLS_409({
    type: "submit",
    color: "primary",
    label: "Salva Azienda",
    ...{ class: "q-px-xl rounded-lg shadow-sm" },
    noCaps: true,
    loading: (__VLS_ctx.savingCompany),
}));
const __VLS_411 = __VLS_410({
    type: "submit",
    color: "primary",
    label: "Salva Azienda",
    ...{ class: "q-px-xl rounded-lg shadow-sm" },
    noCaps: true,
    loading: (__VLS_ctx.savingCompany),
}, ...__VLS_functionalComponentArgsRest(__VLS_410));
/** @type {__VLS_StyleScopedClasses['q-px-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-sm']} */ ;
// @ts-ignore
[vClosePopup, addCompany, companyForm, companyForm, companyForm, companyForm, companyForm, savingCompany,];
var __VLS_374;
var __VLS_375;
// @ts-ignore
[];
var __VLS_368;
// @ts-ignore
[];
var __VLS_356;
// @ts-ignore
[];
var __VLS_350;
// @ts-ignore
[];
var __VLS_3;
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
