_raza_init() {
  if [[ -n "${RAZA_SESSION}" ]]; then
    return
  fi
  RAZA_SESSION=$(raza _shell start-session --host "${HOST}" --user "${USER}")
  echo "CREATED RAZA SESSION $RAZA_SESSION ${RAZA_SESSION}"
  readonly RAZA_SESSION
}

_raza_addhistory() {
  _raza_init
  echo $RAZA_SESSION
  raza _shell add --cmd "$*" --pwd "$(pwd)" --start_time "${EPOCHSECONDS}" --session "$RAZA_SESSION"
}

_raza_precmd() {
  _raza_init
  raza _shell pre --retval "$?" --end_time "${EPOCHSECONDS}" --session "$RAZA_SESSION"
}

add-zsh-hook zshaddhistory _raza_addhistory
add-zsh-hook precmd _raza_precmd