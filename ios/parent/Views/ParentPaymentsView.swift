import SwiftUI

public struct PagoPaNoticeModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let deadlineAndAmount: String
    public let isPaid: BooleanLiteralType

    public init(id: String = UUID().uuidString, title: String, deadlineAndAmount: String, isPaid: Bool = false) {
        self.id = id
        self.title = title
        self.deadlineAndAmount = deadlineAndAmount
        self.isPaid = isPaid
    }
}

public struct ParentPaymentsView: View {
    public var notices: [PagoPaNoticeModel]
    public var onPay: ((String) -> Void)?

    public init(notices: [PagoPaNoticeModel] = [], onPay: ((String) -> Void)? = nil) {
        self.notices = notices
        self.onPay = onPay
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if notices.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(notices) { notice in
                            HStack {
                                VStack(alignment: .leading, spacing: 4) {
                                    Text(notice.title)
                                        .fontWeight(.bold)
                                    Text(notice.deadlineAndAmount)
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                }
                                Spacer()
                                if notice.isPaid {
                                    Text("Pagato")
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                } else {
                                    Button("PagoPA") {
                                        onPay?(notice.id)
                                    }
                                    .buttonStyle(.borderedProminent)
                                    .tint(Color.blue)
                                }
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("parent_dashboard_title", comment: ""))
        }
    }
}
