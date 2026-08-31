import SwiftUI

struct SecretaryTimetableView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Orario Scolastico Istituzionale")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Orario Definitivo 2025/2026")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Pubblicato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("32 Classi • Lunedì - Venerdì (6h/g)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("48 cattedre • Conflitti: 0")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Orario e Cattedre")
        }
    }
}
