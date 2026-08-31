import SwiftUI

struct TeacherRubricsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Griglie e Rubriche di Valutazione")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Griglia Esposizione Orale")
                                .fontWeight(.bold)
                            Spacer()
                            Text("4 Livelli")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Proprietà di linguaggio, logica, critica")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Dipartimento Umanistico • Italiano, Storia")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Rubriche Valutative")
        }
    }
}
