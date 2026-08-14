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
  help: {
    title: "帮助中心",
    subtitle: "使用电子注册表的指南和教程",
    searchPlaceholder: "在指南中搜索...",
    noResults: "未找到相关文章",
    categories: "分类",
    allTopics: "所有主题",
    restartTour: "重新开始引导导览",
    openHelp: "打开指南",
    needHelp: "需要帮助？",
    contactSupport: "联系支持",
    fabTooltip: "帮助和指南",
    new: "新"
  }

,
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
  guideCenter: {
    "title": "Help & Knowledge Center",
    "subtitle": "Comprehensive guides for every section of the electronic register",
    "search": "Search in guides...",
    "noResults": "No guides found for",
    "readingTime": "min read",
    "step": "Step",
    "tip": "Tip",
    "warning": "Warning",
    "shortcut": "Shortcut",
    "teacher": {
        "dashboard": {
            "title": "Teacher Dashboard",
            "desc": "Overview of daily lessons, notifications, and quick actions.",
            "content": "The teacher dashboard provides a complete overview of your day. Top widgets display today’s scheduled classes with times and room numbers.\n\nStep 1: Customize your dashboard widgets by dragging them.\nStep 2: Click on any class in the timetable to open its register directly.\nStep 3: Filter notifications by urgency or category.\n\nTip: Enable priority alerts to stay informed about urgent announcements."
        },
        "attendance": {
            "title": "Attendance Register",
            "desc": "How to record attendance, absences, and delays quickly and accurately.",
            "content": "The attendance register automatically loads the student list for the current period.\n\nStep 1: Click on a student name to toggle status (Present, Absent, Late, Early exit).\nStep 2: Add optional justification notes or remarks.\nStep 3: Click \"Sign Lesson\" to confirm and digitally sign the session.\nStep 4: Parents of absent students are notified automatically.\n\nShortcut: Press Ctrl+Enter to quickly save attendance."
        },
        "grades": {
            "title": "Grade Management",
            "desc": "Enter and manage written, oral, and practical grades.",
            "content": "The gradebook is an interactive grid organized by students and evaluation dates.\n\nStep 1: Select class and subject from the dropdown.\nStep 2: Click on any cell to input or edit a grade.\nStep 3: Choose evaluation type (Written, Oral, Practical) and score.\nStep 4: Add private or public feedback comments.\nStep 5: Subject and class averages update instantly.\n\nWarning: Published grades are visible to students and parents within seconds."
        },
        "agenda": {
            "title": "School Agenda & Homework",
            "desc": "Plan tests, assignments, and shared class activities.",
            "content": "The teacher agenda syncs automatically with students and parents.\n\nStep 1: Click on any calendar date to schedule an event.\nStep 2: Choose event type: Homework, Exam, Field trip, or Reminder.\nStep 3: Select subject, description, and attach learning files.\nStep 4: Check \"Notify students\" to broadcast an instant push notification.\n\nTip: Use the Week view to check for test overlaps across subjects."
        },
        "communications": {
            "title": "Communications & Messages",
            "desc": "Send structured messages to classes, parents, and colleagues.",
            "content": "Send announcements through dedicated channels: direct messages, class circulars, or urgent bulletins.\n\nStep 1: Navigate to \"Communications\" → \"New Message\".\nStep 2: Choose channel type and recipient list.\nStep 3: Format your message using the rich-text editor.\nStep 4: Send immediately or schedule for future publication.\n\nWarning: Urgent announcements trigger push notifications and SMS alerts."
        },
        "meetings": {
            "title": "Parent-Teacher Meetings",
            "desc": "Configure availability slots and manage booking appointments.",
            "content": "Manage weekly office hours and individual parent conferences.\n\nStep 1: Go to \"Meetings\" → \"Set Availability\".\nStep 2: Define weekly recurring time slots and appointment durations.\nStep 3: Parents book online; you receive automated calendar invites.\nStep 4: Review student performance records prior to the meeting.\n\nTip: Enable the integrated video call link for virtual meetings."
        }
    },
    "student": {
        "dashboard": {
            "title": "Student Dashboard",
            "desc": "Your daily hub for deadlines, grades, and schedules.",
            "content": "The student dashboard prioritizes urgent tasks and upcoming exams.\n\nStep 1: Review \"Today’s Schedule\" for lessons and rooms.\nStep 2: Check \"Due Today\" and \"Tomorrow’s Exams\".\nStep 3: Monitor your overall GPA and recent grade additions.\n\nTip: Install the PWA on your phone for instant notification alerts."
        },
        "grades": {
            "title": "My Grades & Performance",
            "desc": "Track your subject grades, averages, and evaluation history.",
            "content": "View all assessments organized by subject and academic term.\n\nStep 1: Select a subject to inspect detailed teacher feedback.\nStep 2: Review the trend line to understand your progress over time.\nStep 3: Switch between terms (1st Term, 2nd Term, Final).\nStep 4: Download a full PDF transcript at any time."
        },
        "homework": {
            "title": "Homework & Assignments",
            "desc": "Organize your homework, downloads, and submission deadlines.",
            "content": "Stay on top of all homework assignments and study materials.\n\nStep 1: Filter tasks by due date or subject.\nStep 2: Click on an assignment to download attachments.\nStep 3: Mark tasks as \"Completed\" to keep track of your workload."
        }
    },
    "parent": {
        "monitoring": {
            "title": "Monitoring Your Children",
            "desc": "Follow academic performance, daily attendance, and notices in real time.",
            "content": "Access comprehensive academic information for all enrolled children.\n\nStep 1: Select your child from the top switcher.\nStep 2: Check daily attendance and real-time grade notifications.\nStep 3: Read teacher comments on individual assessments.\nStep 4: View academic progress charts.\n\nTip: Configure notification thresholds to receive alerts if grades drop."
        },
        "meetings": {
            "title": "Booking Teacher Meetings",
            "desc": "Schedule and manage conference appointments with teachers.",
            "content": "Easily book one-on-one parent-teacher conferences.\n\nStep 1: Open \"Meetings\" and select the desired teacher.\nStep 2: Pick an available green time slot on the calendar.\nStep 3: Confirm booking to receive an email reminder.\n\nWarning: Cancellations should be submitted at least 2 hours in advance."
        }
    },
    "secretary": {
        "students": {
            "title": "Student Registry Management",
            "desc": "Search, manage, and update student profiles and demographic records.",
            "content": "The student registry is the central database for all enrolled learners.\n\nStep 1: Use the global search bar (name, tax code, class).\nStep 2: Click a record to open the full student dossier.\nStep 3: Edit contact information, enrollment status, or guardian links.\nStep 4: Export filtered data to Excel or CSV formats."
        },
        "certificates": {
            "title": "Certificate Generation",
            "desc": "Issue official school documents, attendance proofs, and certificates.",
            "content": "Produce signed school documents in seconds.\n\nStep 1: Search for the student record.\nStep 2: Go to \"Certificates\" → \"Generate\".\nStep 3: Select certificate template (Enrollment, Attendance, Transcripts).\nStep 4: Preview, apply digital stamp, and print or email directly."
        }
    },
    "admin": {
        "monitoring": {
            "title": "System & Health Monitoring",
            "desc": "Monitor platform uptime, API response times, and server performance.",
            "content": "Real-time telemetry and health monitoring across all platform services.\n\nStep 1: Inspect server CPU, RAM, and active session count.\nStep 2: Monitor API throughput and error rates over the last 24 hours.\nStep 3: Configure alert thresholds for proactive notifications.\n\nWarning: Use Maintenance Mode during scheduled platform updates."
        },
        "users": {
            "title": "User & Access Management",
            "desc": "Manage user accounts, RBAC permissions, and authentication policies.",
            "content": "Centralized user provisioning and security management.\n\nStep 1: Search users by role, institution, or active status.\nStep 2: Create single accounts or use the bulk CSV/Excel import tool.\nStep 3: Enforce 2FA/MFA and manage password resets.\nStep 4: Review immutable audit logs for compliance auditing."
        }
    }
}
}
