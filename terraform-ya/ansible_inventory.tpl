---
all:
  hosts:
  children:
    ya:
      hosts:
    %{for num in drop_num~}
      ${drop_ip[num]}:
            ansible_user: ${drop_user}
    %{endfor~}