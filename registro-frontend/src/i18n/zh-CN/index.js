export default {
  common: {
    welcome: '欢迎',
    back: '返回',
    save: '保存',
    cancel: '取消',
    close: '关闭',
    delete: '删除',
    edit: '编辑',
    search: '搜索 (Ctrl+K)...',
    filter: '筛选',
    actions: '操作',
    loading: '正在加载...',
    success: '操作成功完成',
    error: '发生错误',
    systemSettings: '系统设置',
    language: '语言 / Language',
    security: '安全与身份验证',
    notifications: '通知与提醒',
    mainMenu: '主菜单',
    logout: '退出登录'
  },
  login: {
    welcomeBack: '欢迎回来',
    subtitle: '登录以访问电子记分册系统',
    emailLabel: '电子邮箱',
    emailRequired: '电子邮箱不能为空',
    passwordLabel: '密码',
    passwordRequired: '密码不能为空',
    rememberMe: '记住我',
    submit: '登录',
    noAccount: '还没有账号？',
    contactSecretary: '联系教务处',
    contactTitle: '联系教务处',
    contactSubtitle: '请从下拉菜单中选择您的学校，以查看教务处的联系邮箱和电话。',
    selectSchool: '选择您的学校 / 机构',
    noSchoolFound: '未找到学校',
    emailSegreteria: '教务处邮箱',
    phone: '电话',
    sendEmail: '发送邮件',
    copyEmail: '复制邮箱',
    emailCopied: '邮箱地址已复制到剪贴板！',
    chooseSchoolPrompt: '请从上方菜单选择学校以查看联系信息。',
    sessionExpired: '登录已过期，请重新登录。',
    showPassword: '显示密码',
    hidePassword: '隐藏密码'
  },
  nav: {
    dashboard: '仪表盘',
    schools: '学校管理',
    mySchool: '我的学校',
    users: '用户管理',
    admins: '管理员',
    monitoring: '系统监控',
    analytics: '全局分析',
    auditLogs: '审计日志',
    settings: '设置',
    support: '技术支持',
    featureFlags: '学校功能配置',
    elearning: 'Google与Teams网课集成',
    students: '学生管理',
    classes: '班级管理',
    groups: '语言/分层小组',
    documents: '文档',
    certificates: '证明与证书',
    textbooks: '教材管理',
    meetings: '会议管理',
    communications: '通知公告',
    reports: '统计报告',
    pcto: '社会实践 / PCTO',
    scrutiny: '期末考评与审查',
    myClasses: '我的班级',
    classRegister: '课堂日志',
    uda: 'UdA教学计划',
    competencies: '能力评估',
    grades: '成绩管理',
    attendance: '考勤记录',
    didactics: '教学资料',
    pdp: 'PDP / PEI 个性化计划',
    rubrics: '评分标准',
    coordination: '班主任工作',
    timetable: '课程表',
    agenda: '日程与作业',
    colloqui: '家长约见',
    substitutions: '代课安排',
    verbali: '会议纪要',
    notes: '纪律违规记录',
    myGrades: '我的成绩',
    myAttendance: '我的考勤',
    homework: '家庭作业',
    orientamento: '升学指导',
    calendar: '校历',
    reportCard: '成绩单',
    profile: '个人资料',
    myChildren: '我的孩子',
    goals: '目标',
    trips: '外出与游学',
    payments: '缴费管理',
    assemblies: '大会与会议'
  },
  categories: {
    anagraficheClassi: '学籍与班级',
    attiCertificati: '公文与证书',
    serviziReport: '服务与报告',
    didatticaValutazione: '教学与评估',
    organizzazioneOrario: '组织与课表',
    comunicazioniAtti: '通知与文件',
    percorsiComunicazioni: '发展路径与通知',
    valutazioneDidattica: '评估与教学',
    serviziOrari: '服务与时间表',
    comunicazioniAccount: '通知与账户'
  },
  roles: {
    admin: '管理员',
    superadmin: '超级管理员',
    secretary: '教务人员',
    teacher: '教师',
    student: '学生',
    parent: '家长',
    user: '用户'
  },
  notifications: {
    title: '通知与消息',
    logoutSuccess: '成功退出登录',
    logoutError: '退出登录时发生错误',
    settingsSaved: '设置保存成功！',
    emailCopied: '邮箱已复制！',
    passwordUpdated: '密码更新成功！',
    languageChanged: '语言更新成功'
  },
  errors: {
    connectionError: '服务器连接错误，请检查您的网络。',
    forbidden: '您没有权限执行此操作。',
    serverError: '服务器发生错误，请稍后再试。',
    unauthorized: '登录会话无效或已过期。',
    invalidCredentials: '邮箱或密码不正确。',
    rateLimit: '尝试次数过多，请几分钟后再试。',
    accountDisabled: '账号已被禁用或暂停，请联系教务处。',
    userNotFound: '未找到该用户。',
    passwordMismatch: '两次输入的密码不一致。',
    sessionInvalid: '会话无效，请重新登录。',
    ERR_CURRENT_PASSWORD_INCORRECT: '当前密码不正确。',
    ERR_PASSWORD_COMPLEXITY: '密码必须包含至少一个大写字母、一个小写字母、一个数字和一个特殊字符。',
    ERR_PASSWORD_TOO_SHORT: '密码长度至少为10个字符。',
    ERR_PASSWORD_TOO_LONG: '密码过长（最多128个字符）。',
    ERR_PASSWORD_RECENTLY_USED: '新密码不能与最近使用的5个密码相同。',
    ERR_REQUIRED_FIELDS: '请填写所有必填的密码字段。'
  },
  settings: {
    title: '系统设置',
    subtitle: '配置常规偏好、安全与通知选项',
    languageLabel: '语言 / Language',
    securityLabel: '安全与身份验证',
    securitySub: '密码策略、双重验证 (2FA/MFA) 及会话超时设置',
    notificationsLabel: '通知与提醒',
    notificationsSub: '邮件、推送通知及代课提醒偏好',
    mfaActive: '双重验证已启用',
    mfaOptional: '双重验证可选',
    mfaTitle: '强制要求 2FA / MFA',
    mfaCaption: '要求所有教职工和管理员使用 TOTP 进行双重验证。',
    minPasswordLength: '密码最小长度',
    sessionTimeout: '会话空闲超时',
    maxLoginAttempts: '锁定前的最大失败尝试次数'
  },
  search: {
    placeholder: '搜索学生、教师、教务、班级、通知、菜单项…',
    hint: '输入内容以搜索学生、教师、教务、班级、通告或菜单',
    noResults: '未找到结果：',
    navigate: '导航',
    open: '打开',
    close: '关闭'
  }
,
  onboarding: {
    welcomeTitle: "欢迎使用电子注册表！",
    welcomeSubtitle: "通过几个简单步骤，了解如何充分利用您角色的所有功能。",
    startTour: "开始导览",
    skipTour: "跳过导览",
    next: "下一步",
    prev: "上一步",
    finish: "开始使用注册表",
    stepOf: "第 {current} 步，共 {total} 步",
    restartTour: "重新开始导览",
    tourCompleted: "导览完成！",
    tourCompletedMsg: "您已准备好使用电子注册表。您可以随时查看指南。",
    teacher: {
      step1_title: "教师控制台",
      step1_desc: "您的个人控制台显示今日课程摘要、最新通知以及最常用功能的快速访问入口。",
      step2_title: "班级日志与考勤",
      step2_desc: "在\"我的班级\"中，您可以访问日志、记录出勤和缺勤情况、输入课程主题。",
      step3_title: "成绩管理",
      step3_desc: "在\"成绩\"中，您可以输入口头和书面成绩，查看班级平均分和成绩分布。",
      step4_title: "日程与通讯",
      step4_desc: "日程允许规划测试、作业和活动。在\"通讯\"中，您可以向学生和家长发送消息。",
      step5_title: "家长会面",
      step5_desc: "管理与家长的个别会面：查看预约、可用时间和班级进度。",
      step6_title: "设置与个人资料",
      step6_desc: "在\"设置\"中，您可以自定义语言、界面主题、通知并更新密码。"
    },
    student: {
      step1_title: "您的控制台",
      step1_desc: "个人控制台显示待完成的作业、最新通知和学习成绩摘要。",
      step2_title: "我的成绩",
      step2_desc: "在\"我的成绩\"中，您可以查看教师记录的所有成绩、各科平均分和随时间的进步情况。",
      step3_title: "我的考勤",
      step3_desc: "跟踪您的出勤、缺勤和迟到情况。可以查看每日和每月详情。",
      step4_title: "作业与教学",
      step4_desc: "在\"作业\"中找到教师布置的所有作业及截止日期。\"教学\"显示教师上传的材料。",
      step5_title: "成绩单与文件",
      step5_desc: "从\"成绩单\"可以查看评估文件。在\"文件\"中找到通知和学校材料。",
      step6_title: "学校日历",
      step6_desc: "学校日历显示假期、计划的考试日期和重要的学校活动。"
    },
    parent: {
      step1_title: "家长控制台",
      step1_desc: "您的控制台显示孩子的进步摘要、最新学校通知和未读消息。",
      step2_title: "我的孩子",
      step2_desc: "在\"我的孩子\"中找到已注册孩子的列表。选择一个查看完整的学术档案。",
      step3_title: "成绩与考勤",
      step3_desc: "实时监控孩子的成绩和考勤。立即收到缺勤和新成绩的通知。",
      step4_title: "学校通讯",
      step4_desc: "所有官方学校通讯（通知、教师消息）都汇集在这里。",
      step5_title: "教师会面",
      step5_desc: "直接通过应用程序预约与教师的个别会面，查看可用时间并获得自动确认。",
      step6_title: "付款与文件",
      step6_desc: "管理学校付款并访问孩子的文件（成绩单、证书）。"
    },
    secretary: {
      step1_title: "秘书处控制台",
      step1_desc: "控制台显示待处理活动、最新请求和学校主要统计数据。",
      step2_title: "班级与学生管理",
      step2_desc: "在\"班级\"中管理所有学校班级。在\"学生\"中找到完整的注册表和高级搜索。",
      step3_title: "证书与文件",
      step3_desc: "生成并打印注册和出勤证书。管理数字文档档案。",
      step4_title: "时间表管理",
      step4_desc: "配置学校时间表，管理代课教师并规划活动。",
      step5_title: "报告与统计",
      step5_desc: "生成关于出勤、成绩和注册的自定义报告。",
      step6_title: "用户与通讯",
      step6_desc: "管理教师、学生和家长账户。向所有用户类别发送官方通讯。"
    },
    admin: {
      step1_title: "管理员控制台",
      step1_desc: "控制台显示系统状态、最新活动和平台主要使用指标。",
      step2_title: "系统监控",
      step2_desc: "实时监控系统性能、活跃会话、系统日志和服务状态。",
      step3_title: "用户与学校管理",
      step3_desc: "管理所有用户账户，创建新学校，分配角色和权限。",
      step4_title: "分析与报告",
      step4_desc: "查看平台使用的全局分析，生成详细报告并监控趋势。",
      step5_title: "审计日志",
      step5_desc: "访问系统中所有操作的完整记录，用于审计和合规。",
      step6_title: "系统设置",
      step6_desc: "配置安全策略、双因素身份验证和电子学习集成。"
    }
  },
  onboardingExtra: {
    "openGuide": "Open Full Guide",
    "viewAllFeatures": "Discover all features",
    "completionTitle": "You're ready! 🎉",
    "completionDesc": "You've completed the tour. You can access the full guide at any time from the ? button in the top bar.",
    "teacher": {
        "step7_title": "Digital Register & Signatures",
        "step7_desc": "Digitally sign each lesson with a single click. The register automatically tracks lesson topics, teaching hours, and private notes.",
        "step7_bullets": [
            "Digital lesson signature",
            "Lesson topics",
            "Private notes",
            "Hours log"
        ],
        "step8_title": "Competencies & Learning Units",
        "step8_desc": "Assess student competencies according to EU descriptors, manage Learning Units (UDA) and annual work plans.",
        "step8_bullets": [
            "EU Competencies",
            "Learning Units (UDA)",
            "Annual work plan",
            "Evaluation rubrics"
        ],
        "step1_bullets": [
            "Today's lessons",
            "Recent notifications",
            "Quick shortcuts",
            "Class overview"
        ],
        "step2_bullets": [
            "Attendance register",
            "Lesson topics",
            "Digital signature",
            "Absence management"
        ],
        "step3_bullets": [
            "Oral & written grades",
            "Class average",
            "Distribution chart",
            "Export tables"
        ],
        "step4_bullets": [
            "Homework & tests",
            "Activity calendar",
            "Class announcements",
            "Parent notices"
        ],
        "step5_bullets": [
            "Meeting availability",
            "Parent bookings",
            "Video conferences",
            "Meeting history"
        ],
        "step6_bullets": [
            "Theme & language",
            "Push notifications",
            "Password management",
            "Public profile"
        ]
    },
    "student": {
        "step7_title": "Internships & Competencies",
        "step7_desc": "Track your work-based learning hours (PCTO), view competence certifications, and build your digital portfolio.",
        "step7_bullets": [
            "Internship hours",
            "Certifications",
            "Digital portfolio",
            "Tutor feedback"
        ],
        "step8_title": "Announcements & Noticeboard",
        "step8_desc": "Read official school circulars, answer surveys, and view personalized announcements from your teachers.",
        "step8_bullets": [
            "School circulars",
            "Surveys & answers",
            "Personal alerts",
            "Teacher notices"
        ],
        "step1_bullets": [
            "Upcoming deadlines",
            "Latest notifications",
            "GPA summary",
            "Upcoming tests"
        ],
        "step2_bullets": [
            "Grades by subject",
            "Live average",
            "Trend chart",
            "Term comparison"
        ],
        "step3_bullets": [
            "Monthly attendance",
            "Absence counter",
            "Late arrivals & early leaves",
            "Excuses & justifications"
        ],
        "step4_bullets": [
            "Assigned homework",
            "Due dates",
            "Learning materials",
            "Links & resources"
        ],
        "step5_bullets": [
            "Digital report card",
            "PDF download",
            "Certificates",
            "Forms"
        ],
        "step6_bullets": [
            "School holidays",
            "Test dates",
            "Extracurriculars",
            "Field trips & events"
        ]
    },
    "parent": {
        "step7_title": "Absence Justifications & Permissions",
        "step7_desc": "Submit digital justifications for student absences directly from the app, and approve early leave authorizations.",
        "step7_bullets": [
            "Online justifications",
            "Early leave permits",
            "Special activities",
            "Absence history"
        ],
        "step8_title": "Academic Analytics & Progress",
        "step8_desc": "View academic performance charts over time, compare subject averages, and track improvement trends.",
        "step8_bullets": [
            "Progress chart",
            "Subject comparison",
            "Historical trends",
            "Goals achieved"
        ],
        "step1_bullets": [
            "Children overview",
            "Recent alerts",
            "Unread messages",
            "Appointments"
        ],
        "step2_bullets": [
            "Academic profile",
            "Student documents",
            "Teacher contacts",
            "Enrollment info"
        ],
        "step3_bullets": [
            "Real-time grades",
            "Subject GPA",
            "Daily attendance",
            "Automated alerts"
        ],
        "step4_bullets": [
            "School circulars",
            "Teacher notices",
            "Urgent alerts",
            "Digital bulletin board"
        ],
        "step5_bullets": [
            "Book parent meeting",
            "Choose time slot",
            "Email confirmation",
            "Reschedule"
        ],
        "step6_bullets": [
            "Tuition & fees",
            "Payment history",
            "Download receipts",
            "Request documents"
        ]
    },
    "secretary": {
        "step7_title": "Enrollment & Student Registry",
        "step7_desc": "Manage yearly school enrollments, maintain student demographic records, and handle transfers and special education plans.",
        "step7_bullets": [
            "Yearly enrollment",
            "Transfers in/out",
            "Special needs / DSA",
            "Registry archives"
        ],
        "step8_title": "Circulars & School Broadcasts",
        "step8_desc": "Compose and publish digital circulars, manage noticeboards, and send targeted group broadcasts with signature tracking.",
        "step8_bullets": [
            "Compose circulars",
            "Digital signature",
            "Targeted recipients",
            "Publication history"
        ],
        "step1_bullets": [
            "Pending tasks",
            "Incoming requests",
            "School statistics",
            "System alerts"
        ],
        "step2_bullets": [
            "Class list",
            "Teacher assignments",
            "Weekly timetable",
            "Revision history"
        ],
        "step3_bullets": [
            "Enrollment certificates",
            "Attendance records",
            "PDF export",
            "Digital seal"
        ],
        "step4_bullets": [
            "School schedule",
            "Teacher substitutions",
            "Extracurriculars",
            "Staff notices"
        ],
        "step5_bullets": [
            "Attendance reports",
            "Official exports",
            "Statistical charts",
            "Advanced filters"
        ],
        "step6_bullets": [
            "Create accounts",
            "Assign roles",
            "Reset passwords",
            "Access control"
        ]
    },
    "admin": {
        "step7_title": "Schools & Campuses Management",
        "step7_desc": "Create and configure multiple school branches, manage tenant settings, assign campus administrators, and monitor activity.",
        "step7_bullets": [
            "New campuses",
            "Per-school config",
            "Branch admins",
            "Multi-tenant"
        ],
        "step8_title": "Integrations & API Management",
        "step8_desc": "Connect third-party systems (Google Workspace, Microsoft 365, LMS), configure webhooks, and monitor API traffic.",
        "step8_bullets": [
            "Google Workspace",
            "Microsoft 365",
            "Webhook config",
            "API monitoring"
        ],
        "step1_bullets": [
            "System health",
            "Active users",
            "Open sessions",
            "Recent errors"
        ],
        "step2_bullets": [
            "CPU & memory usage",
            "API latency",
            "Live logs",
            "Threshold alerts"
        ],
        "step3_bullets": [
            "Manage users",
            "RBAC roles",
            "Manage schools",
            "Bulk CSV import"
        ],
        "step4_bullets": [
            "Platform metrics",
            "Periodic reports",
            "User growth",
            "Data exports"
        ],
        "step5_bullets": [
            "Immutable audit logs",
            "Security filters",
            "Compliance export",
            "GDPR tools"
        ],
        "step6_bullets": [
            "Password policy",
            "Mandatory 2FA",
            "Feature flags",
            "Maintenance mode"
        ]
    }
},
  help: {
    title: '帮助与知识中心',
    subtitle: '电子记分册使用指南与操作教程',
    searchPlaceholder: '搜索帮助文章...',
    noResults: '未找到相关文章：',
    categories: '分类',
    allTopics: '所有主题',
    restartTour: '重新播放引导',
    openHelp: '打开指南',
    needHelp: '需要帮助？',
    contactSupport: '联系技术支持',
    fabTooltip: '帮助与指南',
    new: '新功能',
    teacher: {
      cat_register: '班级考勤簿',
      cat_grades: '成绩管理',
      cat_attendance: '出勤管理',
      cat_agenda: '日程与布置',
      cat_settings: '设置',
      q1: '如何记录班级出勤与缺勤？',
      a1: '进入“我的班级” → 选择班级 → 点击“出勤”标签页。为每位学生点击 出 (出席)、缺 (缺勤)、迟 (迟到)。最后点击“数字签名”保存。',
      q2: '如何使用键盘快速网格 (Matrix View) 输入成绩？',
      a2: '在“成绩管理”中选择班级与科目，开启 Matrix View 网格模式。使用 TAB 和方向键切换单元格，输入分数后按 ENTER。',
      q3: '如何在日程表中安排考试或布置作业？',
      a3: '在“日程表”中点击日期 → “添加事件”。选择类型（作业、考试、研学活动），选择科目并保存。',
      q4: '如何向家长发送紧急通知或通告？',
      a4: '进入“消息中心” → “新建通知”。选择接收对象，勾选“需要阅读回执”并发送。',
      q5: '如何设置家长会与家长接待时间？',
      a5: '在“预约接待” → “设置可用时段”。设置每周开放的日期、时间段及每次时长（15、20或30分钟）。',
      q6: '如何进入期末评议与终考评定面板 (Scrutinio)？',
      a6: '若您是班主任，侧边栏将显示“期末评议”。输入拟定成绩、品行评语及综合评语。',
      q7: '如何制定或查看个性化教育计划 (PDP / PEI)？',
      a7: '从班级菜单进入“个性化计划 (PDP)”。配置补偿性措施与豁免条款。',
      q8: '如何恢复未保存的考勤簿草稿？',
      a8: '系统每15秒自动保存一次草稿。重新打开页面时，系统会提示您恢复未保存的内容。',
      q9: '如何切换界面语言与暗黑模式？',
      a9: '在用户菜单的“设置”中，可在9种语言间无缝切换，并开启暗黑模式或 OpenDyslexic 易读字体。',
      q10: '如何撤销刚刚误输入的成绩？',
      a10: '保存成绩后，屏幕底部将出现15秒内有效的“撤销”浮条。超时后可直接在单元格修改。'
    },
    student: {
      cat_grades: '我的成绩',
      cat_attendance: '我的出勤',
      cat_homework: '作业列表',
      cat_documents: '文档与成绩单',
      cat_settings: '个人设置',
      q1: '如何查看我的成绩和各科平均分？',
      a1: '在侧边栏打开“我的成绩”，可查看笔试、口试及实操考核的详细记录与加权平均分。',
      q2: '成绩模拟器 (Simulateur) 如何工作？',
      a2: '在“我的成绩”中点击“成绩模拟器”。输入未来考试的预测分数，系统将计算目标平均分。',
      q3: '如何查看我的累计缺勤学时？',
      a3: '在“我的出勤”中，计数器会显示累计缺勤学时以及距离年度上限的剩余额度。',
      q4: '在何处查看老师布置的作业？',
      a4: '在“作业”页面或仪表板。作业按截止日期排序，并附有老师上传的学习资料。',
      q5: '如何下载季度成绩单 PDF？',
      a5: '进入“成绩单与文档”，选择学期并点击“下载成绩单 PDF”。',
      q6: '在何处查看校外实习学时 (PCTO)？',
      a6: '打开“PCTO 实习档案”即可查看认证学时、合作企业及指导老师评语。',
      q7: '如何确认阅读通知或通告？',
      a7: '打开“消息中心”，选择对应通知并点击“确认已读”。',
      q8: '如何开启 OpenDyslexic 读写障碍友好字体？',
      a8: '点击顶部无障碍图标，或进入“个人中心” → “无障碍设置”开启。',
      q9: '忘记密码怎么办？',
      a9: '在登录页面点击“忘记密码？”，输入绑定的电子邮箱重置。',
      q10: '如何使用 Ctrl+K 快捷搜索？',
      a10: '按下 Ctrl+K (Mac 上为 Cmd+K) 打开全局搜索框，可快速跳转至任意功能或课程。'
    },
    parent: {
      cat_monitoring: '子女监控',
      cat_communications: '消息通告',
      cat_meetings: '家长会预约',
      cat_documents: '文档与缴费',
      cat_settings: '设置',
      q1: '如何在多位在读子女间切换？',
      a1: '使用页面顶部的学生切换器，即可在多位子女的档案间瞬间切换。',
      q2: '成绩和缺勤是否有实时推送通知？',
      a2: '是的！每当教师录入新成绩或缺勤记录时，手机应用将收到实时推送和电子邮件。',
      q3: '如何预约与任课老师面谈？',
      a3: '进入“预约接待” → 选择老师 → 在日历中选择绿色可用时段 → 确认预约。',
      q4: '如何通过 PagoPA 缴纳学杂费或活动费？',
      a4: '进入“PagoPA 缴费”。查看待缴账单，支持在线刷卡、PayPal 或下载 QR 码缴费。',
      q5: '如何在在线提交请假条？',
      a5: '在“出勤与请假”中，点击红色缺勤条目旁边的“请假”，选择原因并输入 PIN 码签名。',
      q6: '如何下载孩子的成绩单？',
      a6: '选择孩子 → “成绩单” → 选择学期 → “下载成绩单 PDF”。',
      q7: '未签署的通告在哪里查看？',
      a7: '在“消息中心”。需要家属确认的通知会带有显眼的“签署知悉”按钮。',
      q8: '如何联系学校行政办公室？',
      a8: '在“帮助与支持”中选择“联系客服”，可直接向学校秘书处提交咨询工单。',
      q9: '两位家长可以分别拥有独立的登录账号吗？',
      a9: '可以。学校秘书处可为同一位学生绑定多位法定监护人的独立账号。',
      q10: '在无网络（离线）状态下如何使用？',
      a10: '应用会展示本地缓存数据。离线填写的请假或表单将在网络恢复后自动同步。'
    },
    secretary: {
      cat_students: '学生档案簿',
      cat_classes: '班级编排',
      cat_documents: '证明与公文',
      cat_timetable: '课表与代课',
      cat_reports: '统计与导出',
      q1: '如何快速搜索学生档案？',
      a1: '使用顶部全局搜索栏，支持按姓名、税号、学号或班级快速检索。',
      q2: '如何开具官方在读证明？',
      a2: '在“证明与文档”中搜索学生 → “生成证明 PDF”。生成的公文带有电子印章。',
      q3: '如何处理教师日常临时代课安排？',
      a3: '进入“课表与代课” → 选择日期与缺勤教师。系统会自动推荐无课的空闲教师。',
      q4: '如何导出教育部门 (SIDI) 所需的数据文件？',
      a4: '在“统计与导出” → “教育部 SIDI 导出”中选择时间段与导出格式。',
      q5: '如何为新学生或家长创建系统账号？',
      a5: '在“用户管理” → “新建用户”中填写基本信息并分配角色，系统将发送初始密码邮件。',
      q6: '如何编排新学年的班级与课程？',
      a6: '在“班级管理”中创建新班级，批量导入升学学生并分配任课教师团队。',
      q7: '如何发布学校官方公文通告？',
      a7: '在“消息与通告” → “新建公文”中输入文号、标题、正文，选择目标群体并发布。',
      q8: '如何管理各班级教材选订目录？',
      a8: '在“教材管理”中选择班级与科目，录入 ISBN 书号、书名及出版社信息。',
      q9: '如何核查 PagoPA 学费缴纳状态？',
      a9: '打开“财务与 PagoPA”面板，查看已缴与欠费列表，一键发送催缴提醒。',
      q10: '如何废除误发的发文字号 (Protocollo)？',
      a10: '在公文归档中选择对应文号 → “注明原因并作废”。此操作将记录于 Audit Log 审计日志。'
    },
    admin: {
      cat_monitoring: '系统监控',
      cat_users: '用户管理',
      cat_schools: '校区管理',
      cat_security: '安全与认证',
      cat_analytics: '辍学预警分析',
      q1: '如何监控微服务与数据库运行状态？',
      a1: '在“系统监控”中查看 CPU/内存实时图表、PostgreSQL 连接池状态及 API 响应延迟。',
      q2: '如何添加新校区或分校？',
      a2: '进入“校区管理” → “添加校区”。填写学校代码、名称、地址与官方邮箱。',
      q3: '如何强制重置用户密码或开启双重认证 (2FA)？',
      a3: '在“用户管理”中选择账号 → “安全” → 点击“强制重置密码”或“开启 2FA”。',
      q4: '如何查看不可篡改的系统审计日志 (Audit Log)？',
      a4: '进入“审计与合规”。可按操作人、事件类型、日期及 IP 地址过滤查询历史记录。',
      q5: '如何管理平台功能开关 (Feature Flags)？',
      a5: '在“Feature Flags 管理”中，可针对不同校区开启或关闭特定模块（如 PagoPA、SSO 等）。',
      q6: '如何分析学生辍学与慢性缺勤风险？',
      a6: '打开“Analytics & BI”。算法会自动标记缺勤率超过 20% 或平均分低于 6.0 的高风险学生。',
      q7: '如何配置 Google 或 Microsoft 单点登录 (SSO)？',
      a7: '在“集成与 SSO”中输入 Google Cloud 或 Azure 开发者后台提供的 OAuth2 Client ID 与 Secret。',
      q8: '如何执行数据库手动备份？',
      a8: '在“数据库与备份”中一键触发 PostgreSQL 镜像备份，或配置 AWS S3 自动定时备份。',
      q9: '如何查看全站 RBAC 权限矩阵？',
      a9: '进入“角色与权限”，查看9种系统角色的细粒度权限配置矩阵。',
      q10: '如何将平台切换至维护模式？',
      a10: '在“系统设置”中勾选“维护模式”。开启后，普通用户登录将看到系统维护提示。'
    }
  },
  guideCenter: {
    title: "帮助与知识中心",
    subtitle: "电子记分册各模块的详细操作指南",
    search: "搜索指南文章...",
    noResults: "未找到相关指南：",
    readingTime: "分钟阅读",
    step: "步骤",
    tip: "提示",
    warning: "注意",
    shortcut: "快捷键",
    teacher: {
      dashboard: { title: "教师工作台", desc: "今日课程、通知及快捷操作概览。", content: "教师工作台为您展示全天教学安排。顶部卡片实时显示今日课表、上课时间及教室编号。\n\n步骤 1: 拖拽卡片自定义工作台布局。\n步骤 2: 点击任意课程直接打开考勤簿。\n步骤 3: 按紧急程度筛选通知消息。\n\n提示: 开启高优先级通知可实时接收紧急公告。" },
      attendance: { title: "班级考勤簿", desc: "快速准确地点名出勤、缺勤及迟到。", content: "考勤簿会自动加载当前课时的学生名单。\n\n步骤 1: 点击学生姓名切换状态（出席、缺勤、迟到）。\n步骤 2: 填写备注或请假原因。\n步骤 3: 点击“数字签名”确认并完成课时签署。\n步骤 4: 缺勤学生的家长将自动收到提醒通知。\n\n快捷键: 按 Ctrl+Enter 快速保存出勤记录。" },
      grades: { title: "成绩管理与 Matrix 网格", desc: "使用键盘网格极速录入笔试及口试成绩。", content: "成绩簿提供经典视图与键盘极速网格 Matrix View。\n\n步骤 1: 选择班级与科目。\n步骤 2: 开启 Matrix 视图，使用 TAB 与方向键快速移动。\n步骤 3: 输入成绩（1-10分）并按 ENTER。\n步骤 4: 添加公开或私密评语。\n步骤 5: 科目平均分将实时自动重新计算。\n\n注意: 已发布的成绩将在数秒内同步给学生和家长查看。" },
      agenda: { title: "课程表与作业布置", desc: "规划考试、课后作业及班级活动。", content: "教师日程表与学生及家长端自动同步。\n\n步骤 1: 在日历上点击日期新建事件。\n步骤 2: 选择类型：作业、考试、研学活动。\n步骤 3: 填写科目、说明并上传附件资料。\n步骤 4: 若同班级当日已有2场以上考试，系统将自动预警。\n\n提示: 使用周视图可有效避免考试冲突。" },
      communications: { title: "消息中心与通告", desc: "向班级、家长及同事发送结构化通知。", content: "通过专用通道发布消息：私信、班级通告或紧急公告。\n\n步骤 1: 导航至“消息中心” → “新建消息”。\n步骤 2: 选择通道类型与接收人列表。\n步骤 3: 使用富文本编辑器编写内容。\n步骤 4: 立即发送或定时发布。\n\n注意: 紧急消息将触发推送通知。" },
      meetings: { title: "家长会与接待预约", desc: "配置接待时段与预约管理。", content: "管理每周接待日与一对一家长会。\n\n步骤 1: 进入 '预约接待' → '设置可用时段'。\n步骤 2: 设定开放日期、时间及单次时长（15、20、30分钟）。\n步骤 3: 家长在线预约后您将收到确认通知。\n步骤 4: 面谈前可提前调阅学生学情档案。\n\n提示: 可开启集成式在线视频会议链接。" },
      scrutiny: { title: "期末评议与终考 (Scrutinio)", desc: "拟定成绩、品行评语及评议会议记录。", content: "评议汇总表归集所有任课教师的拟定成绩。\n\n步骤 1: 选择班级与评议学期。\n步骤 2: 核对拟定平均分并输入决议成绩。\n步骤 3: 录入品行得分与综合评语。\n步骤 4: 班主任确认无误后锁定评议表。\n步骤 5: 导出并打印官方 PDF 评议记录。\n\n注意: 班主任锁定后，修改需联系教务处解锁。" },
      pdp: { title: "个性化教育计划 (PDP / PEI)", desc: "为特殊需求学生制定补偿性与豁免措施。", content: "PDP 模块协助教师团队制定个性化教学方案。\n\n步骤 1: 从班级菜单进入 'PDP / PEI'。\n步骤 2: 选择学生，配置补偿工具与豁免条款。\n步骤 3: 保存草稿并发送至家长端完成电子签名。\n步骤 4: 考核时激活的措施将以图标形式进行预警提示。\n\n提示: 使用预设评语库可大幅提升撰写效率。" }
    },
    student: {
      dashboard: { title: "学生个人主页", desc: "每日作业、成绩及课表中心。", content: "学生主页优先展示紧急任务与近期考试。\n\n步骤 1: 在“今日”中查看课程表与教室安排。\n步骤 2: 检查今日到期作业与明日考试。\n步骤 3: 关注您的整体平均分变化。\n\n提示: 推荐在手机上安装 PWA 应用以获取实时提醒。" },
      grades: { title: "我的成绩与表现", desc: "追踪各科成绩、平均分及考核历史。", content: "按科目与学期整理的所有考核成绩。\n\n步骤 1: 点击科目阅读任课教师的评语。\n步骤 2: 查看成绩变化趋势图表。\n步骤 3: 切换不同学期进行对比。\n步骤 4: 随时导出完整 PDF 成绩单。" },
      homework: { title: "作业与任务管理", desc: "管理作业提交、资料下载及截止时间。", content: "全面掌控作业进度与学习资料。\n\n步骤 1: 按截止日期或科目筛选任务。\n步骤 2: 点击任务下载老师上传的附件。\n步骤 3: 完成后勾选“标记为已完成”。" },
      attendance: { title: "出勤与缺勤记录", desc: "统计缺勤学时、迟到记录及请假状态。", content: "查看您的出勤概况与学时额度。\n\n步骤 1: 打开“我的出勤”。\n步骤 2: 核对累计缺勤学时计数器。\n步骤 3: 检查月度日历中的缺勤（红）与迟到（橙）。\n步骤 4: 确认家长是否已在线提交请假条。" },
      documents: { title: "成绩单与官方文档", desc: "查看并下载季度成绩单及学校证明。", content: "获取学校教务处发布的官方带签文档。\n\n步骤 1: 进入“成绩单与文档”。\n步骤 2: 选择学年与学期。\n步骤 3: 点击“在线预览成绩单”。\n步骤 4: 点击“下载 PDF”保存带电子签章的文档。" },
      simulator: { title: "成绩目标模拟器", desc: "计算达到目标平均分所需的最少考试分数。", content: "模拟器为您精准计算下一次考核所需达到的分数。\n\n步骤 1: 选择目标科目。\n步骤 2: 输入期望达到的平均分（如：8.5分）。\n步骤 3: 系统自动计算下次考试需达到的最低分数。" },
      pcto: { title: "PCTO 实习与履历档案", desc: "管理校外实习学时及能力认证。", content: "追踪您的高中校外实习 (PCTO) 累计学时。\n\n步骤 1: 打开“PCTO 实习档案”。\n步骤 2: 查看学时进度条与法定标准对比。\n步骤 3: 查看实习企业信息并下载实习证明。" }
    },
    parent: {
      monitoring: { title: "子女学情实时监控", desc: "实时掌握孩子的成绩、出勤及学校通知。", content: "全面调阅名下所有学生的在校档案。\n\n步骤 1: 在顶部切换选择对应子女。\n步骤 2: 检查今日出勤及最新成绩推送。\n步骤 3: 查看教师在作业及测试中留下的评语。\n步骤 4: 观察学业表现趋势图表。\n\n提示: 可设置预警阈值，当平均分下降时自动接收通知。" },
      meetings: { title: "家长会预约", desc: "在线预约任课教师接待时间。", content: "简单几步完成一对一家长会预约。\n\n步骤 1: 打开“预约接待”并选择教师。\n步骤 2: 在日历中挑选绿色的空闲时段。\n步骤 3: 确认预约，系统将自动发送邮件提醒。\n\n注意: 如需取消预约，请至少提前2小时提交。" },
      communications: { title: "消息通知与公告板", desc: "查阅学校公文、紧急通知及班级消息。", content: "学校所有官方公告集中展示。\n\n步骤 1: 打开“消息中心”。\n步骤 2: 按分类筛选（校办、教师、班级）。\n步骤 3: 阅读正文并下载 PDF 附件。\n步骤 4: 提示需要回执时点击“签署知悉”。" },
      pagopa: { title: "PagoPA 学杂费缴纳", desc: "在线安全缴纳研学活动、餐费及学杂费。", content: "集成 PagoPA 官方安全支付系统。\n\n步骤 1: 进入“PagoPA 缴费”。\n步骤 2: 查看待缴账单清单。\n步骤 3: 点击“立即支付”，支持银行卡或电子支付。\n步骤 4: 或下载带有二维码的缴费凭证至线下支付。" },
      documents: { title: "文档与家长表格", desc: "下载成绩单、证明文件及在线签署授权书。", content: "直接在手机上接收带签名的官方文件。\n\n步骤 1: 进入“文档与成绩单”。\n步骤 2: 下载季度成绩单 PDF 文件。\n步骤 3: 在线填写并数字签名研学活动同意书。" },
      justifications: { title: "在线请假与缺勤说明", desc: "在手机上为孩子的缺勤或迟到提交请假条。", content: "告别纸质假条：使用安全 PIN 码在线请假。\n\n步骤 1: 打开“出勤与请假”。\n步骤 2: 查看红色标记的未请假缺勤记录。\n步骤 3: 点击“请假”，选择原因（生病、私事）并输入 PIN 码确认。" }
    },
    secretary: {
      students: { title: "学生学籍档案管理", desc: "检索、编辑及更新学生基本信息与档案。", content: "学生学籍库是教务管理的核心数据库。\n\n步骤 1: 使用全局搜索栏（姓名、税号、班级）。\n步骤 2: 点击卡片打开学生完整学籍档案。\n步骤 3: 修改联系方式或法定监护人信息。\n步骤 4: 导出筛选后的数据至 Excel 或 CSV。" },
      classes: { title: "班级与教学计划编排", desc: "创建班级、分派课程及指定任课教师。", content: "规划新学年班级与课程分配方案。\n\n步骤 1: 打开“班级管理”。\n步骤 2: 点击“新建班级”添加新班级。\n步骤 3: 绑定教学大纲科目并指定任课教师。\n步骤 4: 指定班主任与教务秘书。" },
      certificates: { title: "开具在读证明与公文", desc: "快速开具带电子印章的官方证明文件。", content: "数秒内自动生成官方签署的学籍证明。\n\n步骤 1: 检索目标学生。\n步骤 2: 进入“证明文件” → “生成”。\n步骤 3: 选择模板（在读证明、出勤证明、成绩单）。\n步骤 4: 预览无误后盖电子公章并打印或发送。" },
      timetable: { title: "课表编排与临时代课", desc: "输入周课表及处理每日教师代课排班。", content: "管理学校总课表及缺勤教师代课安排。\n\n步骤 1: 打开“课表与代课”。\n步骤 2: 在周网格中录入各班级课程。\n步骤 3: 当教师请假时，系统会自动推荐有空闲的教师。\n步骤 4: 确认代课安排并一键通知相关教师。" },
      communications: { title: "学校公文与通告发布", desc: "发布官方公文并追踪签收阅读进度。", content: "向全校教职工、学生及家长发布官方信息。\n\n步骤 1: 打开“消息与通告”。\n步骤 2: 点击“新建公文”。\n步骤 3: 填写发文字号、标题及正文内容。\n步骤 4: 选择接收对象与是否需要阅读签名。" },
      reports: { title: "教育部 (SIDI) 数据报表与导出", desc: "汇总统计数据及国家教育部门标准格式导出。", content: "生成教务统计报表与教育部格式文件。\n\n步骤 1: 打开“统计与报表”。\n步骤 2: 选择报表类型（出勤率、评议结果）。\n步骤 3: 选择导出格式：Excel, CSV, PDF 或 XML SIDI。" }
    },
    admin: {
      monitoring: { title: "系统运行状态监控", desc: "监控平台可用性、API 响应时间及服务器负载。", content: "全站微服务架构的实时遥测监控。\n\n步骤 1: 检查 CPU 使用率、内存及活跃会话数。\n步骤 2: 监控过去24小时的 API 吞吐量与错误率。\n步骤 3: 配置自动化故障预警阈值。\n\n注意: 在计划内更新时请开启维护模式。" },
      users: { title: "用户账号与权限管理", desc: "管理用户账号、RBAC 角色及安全策略。", content: "集中式用户安全与账号全生命周期管理。\n\n步骤 1: 按角色、校区或状态检索用户。\n步骤 2: 单个创建账号或使用 CSV 批量导入。\n步骤 3: 强制开启 2FA/MFA 并管理密码重置。\n步骤 4: 查阅不可篡改的安全审计日志。" },
      schools: { title: "多校区与机构管理", desc: "配置分校区、学校代码及网络参数。", content: "管理多租户 (Multi-tenant) 教育机构架构。\n\n步骤 1: 打开“校区管理”。\n步骤 2: 添加新校区（总校、分校、小学部、中学部）。\n步骤 3: 填写学校官方代码与联系方式。" },
      security: { title: "安全策略与身份认证", desc: "密码复杂度、会话超时、2FA及限流规则。", content: "根据 GDPR 及网络安全规范配置防护策略。\n\n步骤 1: 打开“安全与认证”。\n步骤 2: 设置密码最小长度与定期更换周期。\n步骤 3: 为管理人员强制开启双重认证 (2FA)。" },
      analytics: { title: "大数据分析与辍学预警", desc: "辍学风险与慢性缺勤的预测性分析。", content: "利用商业智能 (BI) 工具精准识别高风险学生。\n\n步骤 1: 打开“Analytics & BI”。\n步骤 2: 查看各班级与科目的缺勤热力图。\n步骤 3: 设置预警触发条件（如：缺勤>20% + 均分<6.0）。" },
      integrations: { title: "教学平台集成与 SSO", desc: "同步 Google Classroom、Microsoft Teams 及 SSO。", content: "将电子记分册无缝对接云端教学平台。\n\n步骤 1: 打开“集成与 SSO”。\n步骤 2: 激活 Google Workspace 或 Microsoft 365 模块。\n步骤 3: 输入 OAuth2 Client ID 与 Secret。" },
      audit: { title: "审计日志与 GDPR 合规", desc: "记录所有访问与修改操作的不可篡改系统日志。", content: "保障数据全生命周期合规与可追溯性。\n\n步骤 1: 打开“审计日志与追溯”。\n步骤 2: 检索历史记录：时间戳、用户、角色、IP、操作内容。\n步骤 3: 导出加密 PDF 审计报告以备检查。" }
    }
  }
}
