# Supabase Cloud Setup Guide

1. **Create Project**
   - Go to [database.new](https://database.new) and create a new project.
   - Note the `Reference ID` and `Database Password`.

2. **Database Schema**
   - Navigate to the **SQL Editor** in the Supabase Dashboard.
   - Copy and paste the contents of the `migrations/` folder in order (001 to 009).
   - Alternatively, use the Supabase CLI:
     ```bash
     supabase link --project-ref <your-project-ref>
     supabase db push
     ```

3. **Storage Buckets**
   - Go to **Storage** in the Dashboard.
   - Create a public bucket named `school-assets` (for logos).
   - Create a private bucket named `documents` (for PDP, PEI types).

4. **Environment Variables**
   - Update `registro-backend/.env` with:
     ```
     SUPABASE_URL=https://<your-project-ref>.supabase.co
     SUPABASE_KEY=<your-anon-or-service-role-key>
     DB_HOST=db.<your-project-ref>.supabase.co
     ```

5. **Auth Configuration**
   - Enable Email/Password provider.
   - Configure SMTP if you want custom emails.
