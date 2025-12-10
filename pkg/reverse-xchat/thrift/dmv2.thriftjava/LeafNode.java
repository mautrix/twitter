package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import android.gov.nist.core.net.C0009a;
import android.gov.nist.javax.sdp.fields.C0015d;
import androidx.work.impl.workers.C9941a;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u00004\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\b\n\u0002\b\u0005\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u000f\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 %2\u00020\u0001:\u0002&%BC\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0007\u001a\u0004\u0018\u00010\u0006\u0012\b\u0010\b\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\t\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\n\u0010\u000bJ\u0017\u0010\u000f\u001a\u00020\u000e2\u0006\u0010\r\u001a\u00020\fH\u0016¢\u0006\u0004\b\u000f\u0010\u0010J\u0012\u0010\u0011\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0011\u0010\u0012J\u0012\u0010\u0013\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0013\u0010\u0012J\u0012\u0010\u0014\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0014\u0010\u0012J\u0012\u0010\u0015\u001a\u0004\u0018\u00010\u0006HÆ\u0003¢\u0006\u0004\b\u0015\u0010\u0016J\u0012\u0010\u0017\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0017\u0010\u0012J\u0012\u0010\u0018\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0018\u0010\u0012JX\u0010\u0019\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0007\u001a\u0004\u0018\u00010\u00062\n\b\u0002\u0010\b\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\t\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u0019\u0010\u001aJ\u0010\u0010\u001b\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u001b\u0010\u0012J\u0010\u0010\u001c\u001a\u00020\u0006HÖ\u0001¢\u0006\u0004\b\u001c\u0010\u001dJ\u001a\u0010!\u001a\u00020 2\b\u0010\u001f\u001a\u0004\u0018\u00010\u001eHÖ\u0003¢\u0006\u0004\b!\u0010\"R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010#R\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010#R\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010#R\u0016\u0010\u0007\u001a\u0004\u0018\u00010\u00068\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0007\u0010$R\u0016\u0010\b\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u0010#R\u0016\u0010\t\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\t\u0010#¨\u0006'"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/LeafNode;", "Lcom/bendb/thrifty/a;", "", "subtree_encryption_public_key", "signature_public_key", "keypair_id", "", "max_supported_protocol_version", "parent_hash", "signature", "<init>", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/Integer;Ljava/lang/String;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "component3", "component4", "()Ljava/lang/Integer;", "component5", "component6", "copy", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/Integer;Ljava/lang/String;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/LeafNode;", "toString", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Ljava/lang/Integer;", "Companion", "LeafNodeAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class LeafNode implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String keypair_id;

    @JvmField
    @InterfaceC88465b
    public final Integer max_supported_protocol_version;

    @JvmField
    @InterfaceC88465b
    public final String parent_hash;

    @JvmField
    @InterfaceC88465b
    public final String signature;

    @JvmField
    @InterfaceC88465b
    public final String signature_public_key;

    @JvmField
    @InterfaceC88465b
    public final String subtree_encryption_public_key;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new LeafNodeAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/LeafNode$LeafNodeAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/LeafNode;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/LeafNode;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/LeafNode;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class LeafNodeAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public LeafNode m83729read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            String string2 = null;
            String string3 = null;
            Integer numValueOf = null;
            String string4 = null;
            String string5 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b != 0) {
                    switch (c11265cMo14127V2.f38393b) {
                        case 1:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string = protocol.readString();
                                break;
                            }
                        case 2:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string2 = protocol.readString();
                                break;
                            }
                        case 3:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string3 = protocol.readString();
                                break;
                            }
                        case 4:
                            if (b != 8) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                numValueOf = Integer.valueOf(protocol.mo14132c4());
                                break;
                            }
                        case 5:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string4 = protocol.readString();
                                break;
                            }
                        case 6:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string5 = protocol.readString();
                                break;
                            }
                        default:
                            C11272a.m14141a(protocol, b);
                            break;
                    }
                } else {
                    return new LeafNode(string, string2, string3, numValueOf, string4, string5);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a LeafNode struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("LeafNode");
            if (struct.subtree_encryption_public_key != null) {
                protocol.mo14136v3("subtree_encryption_public_key", 1, (byte) 11);
                protocol.mo14137w0(struct.subtree_encryption_public_key);
            }
            if (struct.signature_public_key != null) {
                protocol.mo14136v3("signature_public_key", 2, (byte) 11);
                protocol.mo14137w0(struct.signature_public_key);
            }
            if (struct.keypair_id != null) {
                protocol.mo14136v3("keypair_id", 3, (byte) 11);
                protocol.mo14137w0(struct.keypair_id);
            }
            if (struct.max_supported_protocol_version != null) {
                protocol.mo14136v3("max_supported_protocol_version", 4, (byte) 8);
                protocol.mo14122C2(struct.max_supported_protocol_version.intValue());
            }
            if (struct.parent_hash != null) {
                protocol.mo14136v3("parent_hash", 5, (byte) 11);
                protocol.mo14137w0(struct.parent_hash);
            }
            if (struct.signature != null) {
                protocol.mo14136v3("signature", 6, (byte) 11);
                protocol.mo14137w0(struct.signature);
            }
            protocol.mo14134i0();
        }
    }

    public LeafNode(@InterfaceC88465b String str, @InterfaceC88465b String str2, @InterfaceC88465b String str3, @InterfaceC88465b Integer num, @InterfaceC88465b String str4, @InterfaceC88465b String str5) {
        this.subtree_encryption_public_key = str;
        this.signature_public_key = str2;
        this.keypair_id = str3;
        this.max_supported_protocol_version = num;
        this.parent_hash = str4;
        this.signature = str5;
    }

    public static /* synthetic */ LeafNode copy$default(LeafNode leafNode, String str, String str2, String str3, Integer num, String str4, String str5, int i, Object obj) {
        if ((i & 1) != 0) {
            str = leafNode.subtree_encryption_public_key;
        }
        if ((i & 2) != 0) {
            str2 = leafNode.signature_public_key;
        }
        String str6 = str2;
        if ((i & 4) != 0) {
            str3 = leafNode.keypair_id;
        }
        String str7 = str3;
        if ((i & 8) != 0) {
            num = leafNode.max_supported_protocol_version;
        }
        Integer num2 = num;
        if ((i & 16) != 0) {
            str4 = leafNode.parent_hash;
        }
        String str8 = str4;
        if ((i & 32) != 0) {
            str5 = leafNode.signature;
        }
        return leafNode.copy(str, str6, str7, num2, str8, str5);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getSubtree_encryption_public_key() {
        return this.subtree_encryption_public_key;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getSignature_public_key() {
        return this.signature_public_key;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final String getKeypair_id() {
        return this.keypair_id;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final Integer getMax_supported_protocol_version() {
        return this.max_supported_protocol_version;
    }

    @InterfaceC88465b
    /* renamed from: component5, reason: from getter */
    public final String getParent_hash() {
        return this.parent_hash;
    }

    @InterfaceC88465b
    /* renamed from: component6, reason: from getter */
    public final String getSignature() {
        return this.signature;
    }

    @InterfaceC88464a
    public final LeafNode copy(@InterfaceC88465b String subtree_encryption_public_key, @InterfaceC88465b String signature_public_key, @InterfaceC88465b String keypair_id, @InterfaceC88465b Integer max_supported_protocol_version, @InterfaceC88465b String parent_hash, @InterfaceC88465b String signature) {
        return new LeafNode(subtree_encryption_public_key, signature_public_key, keypair_id, max_supported_protocol_version, parent_hash, signature);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof LeafNode)) {
            return false;
        }
        LeafNode leafNode = (LeafNode) other;
        return Intrinsics.m65267c(this.subtree_encryption_public_key, leafNode.subtree_encryption_public_key) && Intrinsics.m65267c(this.signature_public_key, leafNode.signature_public_key) && Intrinsics.m65267c(this.keypair_id, leafNode.keypair_id) && Intrinsics.m65267c(this.max_supported_protocol_version, leafNode.max_supported_protocol_version) && Intrinsics.m65267c(this.parent_hash, leafNode.parent_hash) && Intrinsics.m65267c(this.signature, leafNode.signature);
    }

    public int hashCode() {
        String str = this.subtree_encryption_public_key;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        String str2 = this.signature_public_key;
        int iHashCode2 = (iHashCode + (str2 == null ? 0 : str2.hashCode())) * 31;
        String str3 = this.keypair_id;
        int iHashCode3 = (iHashCode2 + (str3 == null ? 0 : str3.hashCode())) * 31;
        Integer num = this.max_supported_protocol_version;
        int iHashCode4 = (iHashCode3 + (num == null ? 0 : num.hashCode())) * 31;
        String str4 = this.parent_hash;
        int iHashCode5 = (iHashCode4 + (str4 == null ? 0 : str4.hashCode())) * 31;
        String str5 = this.signature;
        return iHashCode5 + (str5 != null ? str5.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.subtree_encryption_public_key;
        String str2 = this.signature_public_key;
        String str3 = this.keypair_id;
        Integer num = this.max_supported_protocol_version;
        String str4 = this.parent_hash;
        String str5 = this.signature;
        StringBuilder sbM11b = C0009a.m11b("LeafNode(subtree_encryption_public_key=", str, ", signature_public_key=", str2, ", keypair_id=");
        C9941a.m12991b(num, str3, ", max_supported_protocol_version=", ", parent_hash=", sbM11b);
        return C0015d.m22a(sbM11b, str4, ", signature=", str5, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}