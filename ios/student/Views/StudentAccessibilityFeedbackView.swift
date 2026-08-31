import SwiftUI

public struct StudentAccessibilityFeedbackView: View {
    @State private var barrierType = "Contrasto / Visibilità"
    @State private var description = ""
    @State private var submitted = false
    public var onSubmitFeedback: ((String, String) -> Void)?

    public init(onSubmitFeedback: ((String, String) -> Void)? = nil) {
        self.onSubmitFeedback = onSubmitFeedback
    }

    public var body: some View {
        NavigationView {
            Form {
                if submitted {
                    Section {
                        VStack(alignment: .leading, spacing: 6) {
                            Text(NSLocalizedString("dashboard_title", comment: ""))
                                .fontWeight(.bold)
                                .foregroundColor(.green)
                        }
                    }
                } else {
                    Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                        Text(NSLocalizedString("grades_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }

                    Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                        TextField(NSLocalizedString("dashboard_title", comment: ""), text: $barrierType)
                        TextEditor(text: $description)
                            .frame(height: 100)

                        Button(NSLocalizedString("dashboard_title", comment: "")) {
                            onSubmitFeedback?(barrierType, description)
                            submitted = true
                        }
                        .disabled(description.isEmpty)
                    }
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
