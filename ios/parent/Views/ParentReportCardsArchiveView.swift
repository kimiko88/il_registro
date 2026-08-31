import SwiftUI

struct ParentReportCardsArchiveView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Schede di Valutazione Ufficiali")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Pagella 1° Quadrimestre")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Sigillo QR")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("A.S. 2025/2026 • Media: 7.8")
                            .font(.subheadline)
                            .foregroundColor(.blue)

                        Button(action: {}) {
                            Label("Scarica PDF con Timbro", systemImage: "arrow.down.doc.fill")
                        }
                        .buttonStyle(.bordered)
                        .padding(.top, 2)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Archivio Pagelle")
        }
    }
}
