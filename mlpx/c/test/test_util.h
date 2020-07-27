#ifndef TEST_UTIL
#define TEST_UTIL

#include<math.h>

#define fail(fmt, ...) do { \
		fprintf(stderr, "TEST FAILED (%s, %s:L%i): ", __func__, __FILE__, __LINE__); \
		fprintf(stderr, fmt, __VA_ARGS__); \
		fprintf(stderr, "\n"); \
		fprintf(stderr, "last error was: %s\n", MLPXGetError()); \
		exit(1); \
	} while(0);

#define str_should_equal(str1, str2) do { \
		if (strcmp((str1), (str2)) != 0) { fail("'%s' ('%s') should equal '%s' ('%s')", \
			str1, #str1, str2, #str2); } \
	} while(0)

#define str_should_not_equal(str1, str2) do { \
		if (strcmp((str1), (str2)) == 0) { fail("'%s' ('%s') should no equal '%s' ('%s')", \
			str1, #str1, str2, #str2); } \
	} while(0)

#define should_equal(v1, v2) do { \
		if ( (v1) != (v2) ) { fail("'%s' should equal '%s'", #v1, #v2); } \
	} while(0)

#define should_equal_epsilon(v1, v2, epsilon) do { \
		if ( fabs(v1 - v2) >= epsilon ) { fail("'%s' should equal '%s' within bound %f", #v1, #v2, epsilon); } \
	} while(0)

#define should_not_equal(v1, v2) do { \
		if ( (v1) == (v2) ) { fail("'%s' should not equal '%s'", #v1, #v2); } \
	} while(0)

#define should_be_true(expr) do { \
		if (!(expr)) { fail("'%s' should have been true and was not", #expr); } \
	} while(0)

#define should_be_false(expr) do { \
		if (!(expr)) { fail("'%s' should have been false and was not", #expr); } \
	} while(0)

#define should_be_null(expr) do { \
		if ((expr) != NULL) { fail("'%s' should have been NULL and was not", #expr); } \
	} while(0)

#define should_not_be_null(expr) do { \
		if ((expr) == NULL) { fail("'%s' should not have been NULL and was not", #expr); } \
	} while(0)

#endif /* TEST_UTIL */
