Procedi con "1. JWT accettato come ?token= nell'URL —

auth/middleware.go:43-46
Il token JWT può essere passato come query parameter. Questo lo espone nei log del server, nel Referer header, nella cache dei proxy e nella history del browser. Un attaccante che legge anche solo i log ottiene sessioni valide.

2. Middleware doppio con chiave contesto sbagliata —

middleware/auth.go
Esiste un middleware inattivo che usa c.Set("userID", ...) (camelCase), mentre il middleware reale usa "user_id" e tutti gli handler leggono "user_id". Se questo middleware venisse usato per errore su una route, tutti gli handler riceverebbero userID vuoto silenziosamente.

🟠 Alti 3. fmt.Printf("DEBUG: ...") su stdout in 8 file —

grades/handler.go
,

auth/middleware.go
, ecc.
Stampano UserID, ParentID, errori JWT su stdout. Su Render/Railway i log della dashboard sono accessibili — informatissimi per un attaccante. Nessun RequireRole su route critiche —

main.go:155-239
DELETE /users/:id/gdpr-delete, POST /users/bulk-delete, PATCH /:id/roles — un qualsiasi utente autenticato (anche uno studente) può tentare di chiamarle. La sicurezza dipende solo dai controlli nel service layer.
