import SwiftUI

struct TeacherPdpView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Piani PDP / PEI della Classe")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("M. R. (PDP - DSA)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Firmato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Compensativi: Mappe concettuali, tempo +30%")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Dispensativi: Dispensa da lettura ad alta voce")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Piani PDP / PEI")
        }
    }
}
