import SwiftUI

struct TeacherAdditionalHoursView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Rendicontazione Ore Eccedenti")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("18 Ore Aggiuntive")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Validate")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("10h Supplenze • 8h Progetto STEM")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Lordo: € 630,00 • Trasmissione NoiPA approvata")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Ore Fondo MOF")
        }
    }
}
