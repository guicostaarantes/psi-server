package appointments_templates

// AppointmentModifiedByPsychologistEmailTemplate is an email template used to tell the patient when a appointment is modified by the psychologist
var AppointmentModifiedByPsychologistEmailTemplate = `<h2>Olá {{ .LikeName }} 😊</h2>
<p>Viemos te informar que uma mudança foi proposta na sua próxima consulta com {{ .PsyFullName }}.</p>
<p>Entre no nosso site para confirmar, alterar ou cancelar essa consulta.</p>
<a href="{{ .SiteURL }}">Ir para o site</a>`
