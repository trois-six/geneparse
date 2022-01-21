#include <stdio.h>

unsigned int bindec(unsigned char *param_1) {
  int local_c;
  unsigned int local_8;
  
  local_8 = 0;
  for (local_c = 0; local_c < 4; local_c = local_c + 1) {
    local_8 = local_8 << 8 | (unsigned int)param_1[local_c];
  }

  return local_8;
}


int main(void) {
    unsigned char nbPersons[4];
    unsigned char sosa[4];
    unsigned char iVar2[4];
    FILE *fp;

    fp = fopen("output/pb_base_info.dat","rb");
    fread(nbPersons, 1, 4, fp);
    fread(sosa, 1, 4, fp);

    printf("nbPersons: 0x%x%x%x%x (hex), %d (dec)\n", nbPersons[0], nbPersons[1], nbPersons[2], nbPersons[3], bindec(nbPersons));
    printf("sosa: 0x%x%x%x%x (hex), %d (dec)\n", sosa[0], sosa[1], sosa[2], sosa[3], bindec(sosa));

    fseek(fp, 0xd, 0);
    fread(iVar2, 1, 4, fp);

    printf("iVar2: 0x%x%x%x%x (hex), %d (dec)\n", iVar2[0], iVar2[1], iVar2[2], iVar2[3], bindec(iVar2));

    fclose(fp);

    return 0;
}