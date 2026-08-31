import SwiftUI

struct StudentPCTOView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Progresso Ore di Alternanza (PCTO)")) {
                    VStack(alignment: .leading, spacing: 8) {
                        HStack {
                            Text("Ore Totali:")
                            Spacer()
                            Text("120 / 150 h")
                                .fontWeight(.bold)
                                .foregroundColor(.purple)
                        }
                        ProgressView(value: 120, total: 150)
                            .tint(.purple)
                    }
                    .padding(.vertical, 4)
                }

                Section(header: Text("Progetti ed Esperienze")) {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("Tech Innovation S.r.l.")
                            .fontWeight(.bold)
                        Text("Tutor Aziendale: Ing. Roberto Neri")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        Text("80 ore completate e validate")
                            .font(.caption2)
                            .foregroundColor(.green)
                    }
                }
            }
            .navigationTitle("Percorsi PCTO")
        }
    }
}
