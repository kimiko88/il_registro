import SwiftUI

struct TeacherUdaView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Unità di Apprendimento (UdA)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("UdA: Cittadinanza Digitale & IA")
                                .fontWeight(.bold)
                            Spacer()
                            Text("In Corso")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Discipline: Informatica, Italiano, Filosofia")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Target: Competenze digitali e pensiero critico")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Programmazione UdA")
        }
    }
}
