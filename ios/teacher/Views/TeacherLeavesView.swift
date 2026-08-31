import SwiftUI

struct TeacherLeavesView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Istanze e Permessi Inoltrati")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Permesso Breve (2 Ore)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Approvato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Venerdì 5 Giugno • 11:00 - 13:00")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Sostituto: Prof.ssa Verdi • Prot: PERM-2026-1049")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Permessi & Ferie")
        }
    }
}
