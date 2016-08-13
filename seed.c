#include <stdio.h>
#include <string.h>
#include <time.h>
#include <random>

const char *g_charSet = "ABCDEFGHJKLMNPQRSTWXYZ01234V6789";
unsigned char g_charLookup[256];

void initCharLookup()
{
	memset(g_charLookup, 0xFF, sizeof(g_charLookup));

	int i=0;
	for(const char *p = g_charSet ; *p ; ++p, ++i)
		g_charLookup[*p] = i;
}

unsigned int string2seed(const std::string &str)
{
	int n = str.size();

	// Must be 9 character long
	if(n != 9) return 0;

	// 5th character must be a whitespace
	if(str[4] != ' ') return 0;

	// Convert into a 64-bit seed
	unsigned char bytes[8];
	for(int i=0 ; i<9 ; ++i)
	{
		if(i == 4) continue;

		int j = i > 4 ? i-1 : i;
		unsigned char b = g_charLookup[(unsigned char)str[i]];

		// Invalid character
		if(b == 0xFF) return 0;

		bytes[j] = b;
	}
	
	// Reduce to 32 bits
	unsigned int r0 = bytes[0];
	r0 = (r0 << 5) | bytes[1];
	r0 = (r0 << 5) | bytes[2];
	r0 = (r0 << 5) | bytes[3];
	r0 = (r0 << 5) | bytes[4];
	r0 = (r0 << 5) | bytes[5];
	r0 = (r0 << 2) | (bytes[6] >> 3);
	r0 ^= 0xFEF7FFD;

	unsigned int r = r0;
	unsigned int d = 0;

	while(r != 0)
	{
		d = (r + d) & 0xFF;
		d = ((d >> 7) + (d << 1)) & 0xFF;
		r >>= 5;
	}

	if(d == (bytes[7] | ((bytes[6] << 5) & 0xFF)))
		return r0;
	else
		return 0;
}

std::string seed2string(unsigned int seed)
{
	// Expand to 64 bits
	unsigned int r = seed;
	unsigned int d = 0;
	while(r != 0)
	{
		d = (r + d) & 0xFF;
		d = ((d >> 7) + (d << 1)) & 0xFF;
		r >>= 5;
	}

	unsigned char bytes[8];

	seed ^= 0xFEF7FFD;
	bytes[0] = (seed >> 27) & 0x1F;
	bytes[1] = (seed >> 22) & 0x1F;
	bytes[2] = (seed >> 17) & 0x1F;
	bytes[3] = (seed >> 12) & 0x1F;
	bytes[4] = (seed >> 7 ) & 0x1F;
	bytes[5] = (seed >> 2 ) & 0x1F;
	bytes[6] = ((d | (seed << 8)) >> 5) & 0x1F;
	bytes[7] = d & 0x1F;

	// Convert into valid seed characters
	std::string out = "0000 0000";

	for(int i=0 ; i<9 ; ++i)
	{
		if(i == 4) continue;
		
		int j = i > 4 ? i-1 : i;
		out[i] = g_charSet[bytes[j]];
	}

	return out;
}

int main(int argc, char *argv[])
{
	initCharLookup();

	std::default_random_engine gen(time(NULL));
	std::uniform_int_distribution<unsigned int> dist;

	unsigned int seed = dist(gen);
	//printf("%d\n", seed);
	//seed = 2;
	printf("%s\n", seed2string(seed).c_str());

	return 0;
}
