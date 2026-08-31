import SwiftUI

struct LoginView: View {
    @Binding var isLoggedIn: Bool
    
    @State private var email = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var errorMessage: String? = nil

    var body: some View {
        ZStack {
            Color(red: 0.12, green: 0.16, blue: 0.23)
                .ignoresSafeArea()
            
            VStack(spacing: 24) {
                VStack(spacing: 12) {
                    Image(systemName: "book.closed.fill")
                        .font(.system(size: 60))
                        .foregroundColor(.blue)
                    
                    Text("Registro Docente")
                        .font(.largeTitle)
                        .fontWeight(.bold)
                        .foregroundColor(.white)
                    
                    Text("Gestione lezioni, presenze e voti")
                        .font(.subheadline)
                        .foregroundColor(.white.opacity(0.8))
                }
                
                VStack(spacing: 16) {
                    HStack {
                        Image(systemName: "envelope.fill")
                            .foregroundColor(.gray)
                        TextField("Email Docente", text: $email)
                            .autocapitalization(.none)
                            .keyboardType(.emailAddress)
                    }
                    .padding()
                    .background(Color(UIColor.systemBackground))
                    .cornerRadius(12)
                    
                    HStack {
                        Image(systemName: "lock.fill")
                            .foregroundColor(.gray)
                        SecureField("Password", text: $password)
                    }
                    .padding()
                    .background(Color(UIColor.systemBackground))
                    .cornerRadius(12)
                    
                    if let errorMessage = errorMessage {
                        Text(errorMessage)
                            .font(.caption)
                            .foregroundColor(.red)
                    }
                    
                    Button(action: {
                        if email.isEmpty || password.isEmpty {
                            errorMessage = "Inserisci email e password"
                            return
                        }
                        isLoading = true
                        isLoggedIn = true
                    }) {
                        if isLoading {
                            ProgressView()
                                .tint(.white)
                        } else {
                            Text("ACCEDI AL REGISTRO")
                                .font(.headline)
                                .fontWeight(.bold)
                                .foregroundColor(.white)
                                .frame(maxWidth: .infinity)
                                .padding()
                                .background(Color.blue)
                                .cornerRadius(12)
                        }
                    }
                }
                .padding(24)
                .background(Color(UIColor.secondarySystemGroupedBackground))
                .cornerRadius(24)
                .padding(.horizontal, 20)
            }
        }
    }
}
