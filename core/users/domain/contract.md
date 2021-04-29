# Contract

April 25, 2021

## Teams

<table>
<tr>
<th>Teams</th>
<th>Memberships</th>
<th>Invites</th>
<th>Projects</th>
</tr>
<tr><td>

| Column    | Type  |
|---------|---------|
| ID      | string  |
| Name    | string  |
| UserID    | string  |
| Created | time    |

</td><td>

| Column | Type         |
|-------------|---------|
| ID          | string  |
| TeamID      | string  |
| UserID      | string  |
| Role        | string  |
| Created     | time    |

</td><td>

| Column      | Type      |
|--------------|----------|
|  ID          | string   |
|  Name        | string   |
|  TeamID      | string   |
|  UserID      | string   |
|  Read        | bool     |
|  Expiration  | time     |
|  Accepted    | bool     |
|  Created     | time     |

</td><td>

| Column       | Type     |
|--------------|----------|
|  ID          | string   |
|  Name        | string   |
|  TeamID      | string   |
|  UserID      | string   |
|  Active      | bool     |
|  Public      | bool     |
|  ColumnOrder | []string |
|  Created     | time     |

</td></tr> </table>
