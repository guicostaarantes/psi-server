package users_templates

// InvitePsychologistEmailTemplate is an email template used when a user is invited to join the platform as a psychologist
var InvitePsychologistEmailTemplate = `<h2>Olá 😊</h2>
<p>Você foi convidado a ingressar no PSI como psicólogo.</p>
<p>Caso queira aceitar o convite, clique no botão abaixo e crie uma senha pra você acessar nosso site.</p>
<a href="{{ .SiteURL }}/cadastro?token={{ .Token }}">Criar uma senha</a>`
